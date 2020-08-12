package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"Ripper/encodekit"
	"Ripper/examples/p2p/protocol"
	"Ripper/mapkit"
	mrand "math/rand"

	golog "github.com/ipfs/go-log"
	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	//gologging "github.com/whyrusleeping/go-logging"
)

// Blockchain is a series of validated Blocks
var (
	Blockchain [5000]*protocol.Node
	// head *protocol.Node
	// tail *protocol.Node
	index_tail int

	mutex = &sync.Mutex{}

	updateMsgQueue = make(chan protocol.BlockUpdate)

	LocalStreamList = mapkit.NewConcurrentSyncMap(64)
	AllStreamList   = mapkit.NewConcurrentSyncMap(64)
)

const (
	FILE_MAX_SIZE int64 = 67108864
)

func makeBasicHost(listenPort int, secio bool, randseed int64) (host.Host, error) {

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	// Generate a key pair for this host. We will use it
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(priv),
	}

	if !secio {
		opts = append(opts, libp2p.EnableRelay())
	}

	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := basicHost.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)
	log.Printf("I am %s\n", fullAddr)
	if secio {
		log.Printf("Now run \"go run main.go -l %d -d %s -secio\" on a different terminal\n", listenPort+1, fullAddr)
	} else {
		log.Printf("Now run \"go run main.go -l %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)
	}

	return basicHost, nil
}

func isBlockValid(newBlock, oldBlock *protocol.Node) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// SHA256 hashing
func calculateHash(block *protocol.Node) string {
	record := strconv.Itoa(int(block.Index)) + block.Timestamp + block.StreamId + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// create a new block using previous block's hash
func generateBlock(tailBlock *protocol.Node, streamId string) *protocol.Node {

	var newBlock *protocol.Node
	t := time.Now()

	newBlock.Index = tailBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.StreamId = streamId
	newBlock.PrevHash = tailBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock
}

func main() {

	t := time.Now()
	genesisBlock := protocol.Node{}
	genesisBlock = protocol.Node{0, t.String(), "", "", calculateHash(&genesisBlock), nil}
	//Blockchain = append(Blockchain, genesisBlock)

	golog.SetAllLoggers(golog.LevelInfo) // Change to DEBUG for extra info
	// Parse options from the command line
	listenF := flag.Int("l", 0, "wait for incoming connections")
	target := flag.String("d", "", "target peer to dial")
	secio := flag.Bool("secio", false, "enable secio")
	seed := flag.Int64("seed", 0, "set random seed for id generation")
	flag.Parse()

	if *listenF == 0 {
		log.Fatal("Please provide a port to bind on with -l")
	}

	// Make a host that listens on the given multiaddress
	ha, err := makeBasicHost(*listenF, *secio, *seed)
	if err != nil {
		log.Fatal(err)
		return
	}
	ha.SetStreamHandler("/p2p/1.0.0", handleStream)
	if *target == "" {
		log.Println("listening for connections")
		// Set a stream handler on host A. /p2p/1.0.0 is
		// a user-defined protocol name.
		// hang forever
		select {}
		/**** This is where the listener code ends ****/
	} else {
		// The following code extracts target's peer ID from the
		// given multiaddress
		ipfsaddr, err := ma.NewMultiaddr(*target)
		if err != nil {
			log.Fatalln(err)
			return
		}

		pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			log.Fatalln(err)
			return
		}

		peerid, err := peer.IDB58Decode(pid)
		if err != nil {
			log.Fatalln(err)
			return
		}
		// Decapsulate the /ipfs/<peerID> part from the target
		// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
		targetPeerAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

		// info, err := peer.AddrInfoFromP2pAddr(maddr)
		// if err != nil {
		// 	log.Fatalln(err)
		// }

		// We have a peer ID and a targetAddr so we add it to the peerstore
		// so LibP2P knows how to contact it

		ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

		log.Println("opening stream")
		// make a new stream from host B to host A
		// it should be handled on host A by the handler we set above because
		// we use the same /p2p/1.0.0 protocol

		s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
		if err != nil {
			log.Fatalln(err)
			return
		}

		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		go writeData(rw)
		go readData(rw)
		select {} // hang forever
	}
}
func handleStream(s net.Stream) {
	log.Println("Got a new stream!")
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go readData(rw)
	go writeData(rw)
}

func writeData(rw *bufio.ReadWriter) {

	go func() {
		for {
			time.Sleep(time.Second)
			mutex.Lock()
			 v, err := json.Marshal(Blockchain)
			if err != nil {
				log.Println(err)
			}
			mutex.Unlock()

			mutex.Lock()
			rw.Write(v)
			rw.Flush()
			mutex.Unlock()

		}
	}()
	cfp, err := os.OpenFile("conf.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer cfp.Close()
	for {
		//将 “本地缓存更改操作” 信息存放在已有列表
		select {
		case buf := <-updateMsgQueue:
			switch buf.Cmd {
			case 50: //本地操作 add
				//将本地的新增资源同步到其他peer
				newBlock := generateBlock(Blockchain[index_tail], buf.StreamId)
				
				if isBlockValid(newBlock, Blockchain[index_tail]) {
					mutex.Lock()
					defer mutex.Unlock()
					if index_tail == 5000 {
						list := make([]protocol.Node, 0)
						//超出主链长度则持久化存储
						for index_tail != 0 {
							list = append(list, *Blockchain[index_tail])
							Blockchain[index_tail] = nil
							index_tail--
						}
						body, err := json.Marshal(list)
						if err != nil {
							log.Fatal(err)
							return
						}
						
						
						
					
						Blockchain[0] = newBlock
					} else {
						index_tail++
						Blockchain[index_tail] = newBlock
					}
				}
				mutex.Lock()
				rw.Write([]byte{})
				rw.Flush()
				mutex.Unlock()

				//将 新增 信息存放在已有列表
				LocalStreamList.Set(buf.StreamId, buf)

			case 40: //处理远端请求




			}
		}

	}
}

func PersistentStorage()  {
	if fileInfo.Size() < FILE_MAX_SIZE{
		fp, err := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer fp.Close()

		_, err = fp.Write(body)
		if err != nil {
			log.Fatal(err)
			return
		}
		fileInfo, _ := os.Stat("data.json")
	} else {

	}
}
func readData(rw *bufio.ReadWriter) {
	for {
		buf := make([]byte, 1024)
		tatol, err := rw.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		if tatol >= 0 {
			mutex.Lock()
			//筛选出 hash和 自己的哈希值近的peer存储操作
			index := encodekit.BytesToInt32(buf[1:5])
			if index > Blockchain[index_tail].Index {

				//更新到主链之后
				body := protocol.PeerUpdateModels{}
				err := encodekit.BinaryRead(buf[:tatol], body)
				if err != nil {
					log.Panic(err)
					continue
				}

				// tail.Next = &protocol.Node{
				// 	Index:    index,
				// 	StreamId: string(body.StreamId[:]),
				// }
				// //更改尾结点
				// tail.Next = tail.Next.Next
			}
			mutex.Unlock()
		}
	}
}

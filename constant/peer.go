package constant

import (

	"os"
	"path/filepath"
	"math/rand"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"

	"Ripper/models"
	"Ripper/retrieve"

	"tools/iokit"
	"tools/mapkit"
)

var (
	PeerList   = mapkit.NewConcurrentSyncMap(64)
	LocalID, LocalRT = InitRoutingTable()
)

func InitRoutingTable() (peer.ID, *retrieve.RoutingTable) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	confDir := dir + PEER_INFO_STORE_PATH
	body, err := iokit.ReadAll(confDir)
	if err != nil {
		panic(err)
	}
	loaclRoutingTable := models.FindNodeResponse{}
	if len(body) == 0 {
		local_id := make([]byte, 32)
		rand.Seed(time.Now().UnixNano())
		rand.Read(local_id)
		loaclRoutingTable.PeerId = local_id
		body, err = proto.Marshal(&loaclRoutingTable)
		if err != nil {
			panic(err)
		}
		iokit.Write(confDir, body)
		loaclRoutingTable.PeerId = local_id
	} else {
		err = proto.Unmarshal(body, &loaclRoutingTable)
		if err != nil {
			panic(err)
		}
	}

	ps := pstore.NewMetrics()
	local_id := retrieve.ConvertPeerID(peer.ID(loaclRoutingTable.PeerId))
	rt, err := retrieve.NewRoutingTable(BUCKETSIZE, local_id, time.Minute, ps, 100*time.Minute, nil)
	if err != nil {
		panic(err)
	}

	nowTim := time.Now().Unix()
	for _, v := range loaclRoutingTable.Peerlist {
		rt.TryAddPeer(peer.ID(v.PeerId),true,true)
		v.TimeStamp = nowTim
		PeerList.Set(string(v.PeerId),v)
	}

	return peer.ID(loaclRoutingTable.PeerId), rt
}

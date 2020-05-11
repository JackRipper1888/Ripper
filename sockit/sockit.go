package sockit

import (
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"strings"
)

const (
	UDP = "udp"
	TCP = "tcp"
)

type Net interface {
	Listen(string)
	Write([]byte, string, string) error
	Read(int) <-chan Peer
}
type udpConn struct {
	//iph *ipv4.Header
	conn         *ipv4.RawConn
	resultStream chan Peer
}

// peer指令
type Peer struct {
	Addr net.UDPAddr
	Data []byte
}

func NewNet(netType string) Net {
	switch netType {
	case UDP:
		udp := new(udpConn)
		udp.resultStream = make(chan Peer, 1)
		return udp
	case TCP:
		//udp := new(tcpConn)
		//return udp
	}
	return nil
}

func checkSum(msg []byte) uint16 {
	sum := 0
	for n := 1; n < len(msg)-1; n += 2 {
		sum += int(msg[n])*256 + int(msg[n+1])
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	var ans = uint16(^sum)
	return ans
}
func (this *udpConn) Listen(address string) {
	Listener, err := net.ListenPacket(UDP, address)
	if err != nil {
		log.Fatal(err)
	}
	defer Listener.Close()
	r, err := ipv4.NewRawConn(Listener)
	if err != nil {
		log.Fatal(err)
	}
	this.conn = r
}

func (this *udpConn) Read(size int) <-chan Peer {
	//var	buff = make([]byte,size)
	//h,_,cmd,_ :=this.conn.ReadFrom(buff)
	//dst := h.Dst
	//peerData := Peer{
	//
	//}
	//this.resultStream <- peerData
	return this.resultStream
}

/**
buff:dup数据body
dst:源ip地址+端口
src:目的ip地址+端口
*/

func (this *udpConn) Write(buff []byte, dst, src string) error {
	dstAddr := strings.Split(dst, ":")
	dstPort := dstAddr[1]
	dstIps := strings.Split(dstAddr[0], ".")

	srcAddr := strings.Split(src, ":")
	srcPort := srcAddr[1]
	srcIps := strings.Split(srcAddr[0], ".")
	//udp伪首部
	var udph = make([]byte, 20)
	//源端口号
	udph[12], udph[13] = byte(dstPort[0]), dstPort[1]
	//目的端口号
	udph[14], udph[15] = byte(srcPort[0]), srcPort[1]
	//udp头长度
	udph[16], udph[17] = 0x00, byte(len(buff)+8)
	//校验和
	udph[18], udph[19] = 0x00, 0x00
	//计算校验值
	check := checkSum(append(udph, buff...))
	udph[18], udph[19] = byte(check>>8&255), byte(check&255)

	//填充ip首部
	iph := &ipv4.Header{
		Version: ipv4.Version,
		//IP头长一般是20
		Len:     ipv4.HeaderLen,
		TOS:     0x00,
		TTL:     64,
		Flags:   ipv4.DontFragment,
		FragOff: 0,
		//17为udp传输格式
		Protocol: 17,
		Checksum: 0,
		TotalLen: ipv4.HeaderLen + len(buff),
		//目的IP
		Dst: net.IPv4([]byte(dstIps[0])[0], []byte(dstIps[1])[0], []byte(dstIps[0])[0], []byte(dstIps[1])[0]),
		//源IP
		Src: net.IPv4([]byte(srcIps[0])[0], []byte(srcIps[1])[0], []byte(srcIps[0])[0], []byte(srcIps[1])[0]),
	}
	//buff为数据
	if err := this.conn.WriteTo(iph, append(udph[12:20], buff...), nil); err != nil {
		log.Fatal(err)
		return nil
	}
	return nil
}

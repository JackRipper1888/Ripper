package task

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"sync"

	quic "github.com/lucas-clemente/quic-go"
)

var (
	// 全局udp conn
	Conn *net.UDPConn

	// 全局queue用于发送到controller的队列
	controllerList = make(chan interface{}, 100)

	bufPool sync.Pool
)

// 监听udp端口
func ListenTask(addr string) {
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	sess, err := listener.Accept(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < 10; i++ {
		dealCmdTask(sess)
	}
}

type Bufbody struct {
	Total int
	Buf   [1024]byte
}

func dealCmdTask(sess quic.Session) {
	stream, err := sess.AcceptStream(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		var buf Bufbody
		body := bufPool.Get()
		if body == nil {
			buf = Bufbody{}
		} else {
			buf = body.(Bufbody)
		}
		buf.Total,err = stream.Read(buf.Buf[:])
		if err != nil {
			fmt.Println(err)
			return
		}
		
		fmt.Println(buf.Buf[:buf.Total])
		stream.Write(buf.Buf[:buf.Total])
	}
}


type loggingWriter struct{ io.Writer }

func (w loggingWriter) Write(b []byte) (int, error) {
	fmt.Printf("Server: Got '%s'\n", string(b))
	return w.Writer.Write(b)
}
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}

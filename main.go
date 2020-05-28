package main

import (
	"flag"
	"fmt"
	"github.com/Ripper/udpkit"
	"github.com/panjf2000/gnet"
	"log"
	_ "net/http/pprof"
	"sync"
	"time"
)

var (
	wg         = sync.WaitGroup{}
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file.")
)

type echoServer struct {
	*gnet.EventServer
	//	pool *goroutine.Pool
}

func (es *echoServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	fmt.Println(srv.ReusePort, srv.CountConnections())
	log.Printf("UDP Echo server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

//var (
//	pool,_ = ants.NewPool(256 * 1024, ants.WithNonblocking(true))
//)

func (es *echoServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	// Use ants pool to unblock the event-loop.
	//_ = es.pool.Submit(func() {
	//	c.SendTo(frame)
	//})
	out = frame
	return
}

//func (es *echoServer) OnShutdown(svr gnet.Server) {
//	fmt.Println("撤退")
//}

func main() {
	var port int = 8091
	//go Peerclient(port)
	//echo := new(echoServer)
	//log.Fatal(gnet.Serve(echo, fmt.Sprintf("tcp://:%d", 9000), gnet.WithMulticore(true)))

	//p := goroutine.Default()

	//poolSize := 256 * 1024
	//pool, _ := ants.NewPool(poolSize, ants.WithNonblocking(true))
	//defer pool.Release()

	go udpkit.Peerclient("192.168.100.200", port)
	time.Sleep(24 * time.Hour)
	//p := goroutine.Default()
	//defer p.Release()

	//echo := &echoServer{}
	////events.pool = p
	//log.Fatal(gnet.Serve(echo,
	//	fmt.Sprintf("udp://127.0.0.1:%d", port),
	//	gnet.WithMulticore(true),
	//	//gnet.WithReusePort(true),
	//	gnet.WithNumEventLoop(runtime.NumCPU()),
	//	))
}

//
//func main() {
//	//flag.Parse()
//	//if *cpuprofile != "" {
//	//	f, err := os.Create(*cpuprofile)
//	//	if err != nil {
//	//	}
//	//	pprof.StartCPUProfile(f)
//	//	defer pprof.StopCPUProfile()
//	//}
//	wg.Add(1)
//	go func() {
//		log.Println(http.ListenAndServe("localhost:7777", nil))
//	}()
//	go peerclient()
//
//	//开启udp服务端监听
//	//go udpkit.ListenUdpTask()
//
//	//开启udp服务端处理池
//	//for i := 0; i<runtime.NumCPU();i++ {
//	//	go udpkit.Worker(i)
//	//}
//
//	//wg.Wait()
//	//ctxkit.CancelAll()
//	//logs.Info("peerclient() start")
//
//
//	//var animal AnimalIF
//	//animal = Factory("cat")
//	//count := 0
//	//for {
//	//	count++
//	//	Get(
//	//		"http://183.60.143.82:3030/peer/select/tracker","0001","0008",map[string]string{
//	//		"streamid":"us",
//	//	})
//	//	if count == 10000{
//	//		return
//	//	}
//	//}
//	//animal.Sleep()
//	//showAnimal(animal)
//	//t2,_ := time.Parse("2016-01-02 15:05:05", "2018-04-23 00:00:06")
//	//t1 := time.Date(2018, 1, 2, 15, 5, 0, 0, time.Local)
//	//t2 := time.Date(2018, 1, 2, 15, 0, 0, 0, time.Local)
//	//
//	//logs.Error(t1.Unix()%(24*3600)%300,t2.Unix()%300)
//}

/*type Score struct {
	Num int
}
*/

package udpkit

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"net"
	"os"
	"sync"
	"time"
)

var (
	Conn *net.UDPConn
	//repPeerList = make(chan *PeerCmd)
	reqPeerList = make(chan *PeerCmd)

	//get,give = bufkit.MakeRecycler()
)

// peer指令
type PeerCmd struct {
	peerAddr *net.UDPAddr
	data     []byte
}

// 监听udp端口
func ListenUdpTask() {
	logPre := "|listen|ips=%s|msg:%s"
	listenIps := "127.0.0.1:8091"
	netAddr, err := net.ResolveUDPAddr("udp", listenIps)
	if err != nil {
		fmt.Println("@"+logPre, listenIps, err.Error())
		//PressEnterToExit()
	}
	Conn, err = net.ListenUDP("udp", netAddr)
	if err != nil {
		fmt.Println("@"+logPre, listenIps, err.Error())
		//PressEnterToExit()
		return
	}
	defer Conn.Close()
	fmt.Println("Task:ListenUdpTask() start")

	num := 100 * 100 * 20
	p := NewWorkerPool(num)
	p.Run()

	var peerAddr *net.UDPAddr
	var countTotal int
	// 插入队列
	var Cmd PeerCmd
	for {
		data := make([]byte, 1024)
		// Here must use make and give the lenth of buffer
		countTotal, peerAddr, err = Conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println("@"+logPre, listenIps, "ReadFromUDP:"+err.Error())
			continue
		}

		Cmd.peerAddr = peerAddr
		Cmd.data = data[:countTotal]
		//Conn.WriteToUDP(Cmd.data, &Cmd.peerAddr)
		// give <- data
		//Conn.WriteToUDP(Cmd.data, &Cmd.peerAddr)
		p.JobQueue <- &Cmd
	}
}

//func WokerPool()  {
//	num := runtime.NumCPU()
//	fmt.Println("CPUNUM:",num)
//	for i:= 0; i<10 ;i++{
//		go	Worker(i)
//	}
//	go addWorder()
//}

//func addWorder()  {
//	for{
//
//	}
//}
//func Worker(i int) {
//	var data *PeerCmd
//	for {
//		select {
//		case data = <- reqPeerList:
//			//repPeerList <- data
//			fmt.Println(i,string(data.data))
//			Conn.WriteToUDP(data.data, data.peerAddr)
//		//case <-time.After(1 bee* time.Minute):
//		//	return
//		}
//	}
//}

//// 返回peer指令
//func RepPeerCmdTask() {
//	logs.Info("Task:RepPeerCmdTask() start")
//	//ctx, _ := ctxkit.CtxAdd()
//	var data PeerCmd
//	for {
//		select {
//		case data = <-repPeerList:
//			Conn.WriteToUDP(data.data, &data.peerAddr)
//		//case <-ctx.Done():
//		//	return
//		}
//	}
//}

func UdpClent(addr string, data map[string]interface{}) ([]byte, error) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		os.Exit(1)
		return nil, err
	}
	defer conn.Close()

	body, _ := json.Marshal(data)
	_, err = conn.Write([]byte(body))
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))
	var msg [1024]byte
	//l, err := conn.Read(msg[0:])
	//if err != nil {
	//	fmt.Println(err)
	//	return nil, err
	//}
	return msg[:], nil
}

//type Job interface {
//	Do()
//}

func (this *PeerCmd) Do() {
	fmt.Println("num:", this.data)
	//time.Sleep(1 * 1 * time.Second)
	Conn.WriteToUDP(this.data, this.peerAddr)
}

type Worker struct {
	JobQueue chan *PeerCmd
}

func NewWorker() Worker {
	return Worker{JobQueue: make(chan *PeerCmd)}
}

func (w Worker) Run(wq chan chan *PeerCmd) {
	go func() {
		for {
			wq <- w.JobQueue
			select {
			case job := <-w.JobQueue:
				job.Do()
			}
		}
	}()

}

type WorkerPool struct {
	workerlen   int
	JobQueue    chan *PeerCmd
	WorkerQueue chan chan *PeerCmd
}

func NewWorkerPool(workerlen int) *WorkerPool {
	return &WorkerPool{
		workerlen:   workerlen,
		JobQueue:    make(chan *PeerCmd),
		WorkerQueue: make(chan chan *PeerCmd, workerlen),
	}
}
func (wp *WorkerPool) Run() {
	fmt.Println("初始化worker")
	//初始化worker
	for i := 0; i < wp.workerlen; i++ {
		worker := NewWorker()
		worker.Run(wp.WorkerQueue)
	}
	// 循环获取可用的worker,往worker中写job
	go func() {
		for {
			select {
			case job := <-wp.JobQueue:
				worker := <-wp.WorkerQueue
				worker <- job
			}
		}
	}()
}

var (
	login = map[string]interface{}{
		"cmd":           "login",
		"appid":         "250",
		"intranet_addr": "192.168.0.11:2020",
		"lanips":        []string{"192.168.0.11:2020"},
		"nattype":       1,
	}
	keepalive = map[string]interface{}{
		"cmd":           "keepalive",
		"appid":         "250",
		"intranet_addr": "192.168.0.11:2020",
		"streamid":      "05071246553639577026050703298691",
		"order":         1,
	}
	peerlist = map[string]interface{}{
		"cmd":      "peerlist",
		"appid":    "250",
		"streamid": "05071246553639577026050703298691",
		"order":    1,
	}
	cache1 = map[string]interface{}{
		"cmd":           "cache",
		"appid":         "250",
		"cache_type":    10,
		"intranet_addr": "192.168.0.11:2020",
		"data": []map[string]interface{}{
			{
				"streamid": "98691",
				"order":    1,
				"lasttime": 1590030441,
			},
			{
				"streamid": "98682",
				"order":    1,
				"lasttime": 1590030441,
			},
		},
		"nattype": 1,
	}
	cache2 = map[string]interface{}{
		"cmd":           "cache",
		"appid":         "250",
		"cache_type":    10,
		"intranet_addr": "192.168.0.11:2020",
		"data": []map[string]interface{}{
			{
				"streamid": "98691",
				"order":    1,
				"lasttime": 1590030441,
			},
			{
				"streamid": "98673",
				"order":    1,
				"lasttime": 1590030441,
			},
		},
		"nattype": 1,
	}
	delete_cache = map[string]interface{}{
		"cmd":           "delete_cache",
		"appid":         "250",
		"intranet_addr": "192.168.0.11:2020",
		"stream_ids":    []string{"98691"},
	}
	state1 = map[string]interface{}{
		"cmd":      "state",
		"appid":    "250",
		"platform": 1,
		"data": map[string]int64{
			"q_recv_source_bytes": 100,
			"q_send_player_bytes": 100,
		},
	}
	state2 = map[string]interface{}{
		"cmd":      "state",
		"appid":    "251",
		"platform": 1,
		"data": map[string]int64{
			"q_recv_source_bytes": 100,
			"q_send_player_bytes": 100,
		},
	}
	state3 = map[string]interface{}{
		"cmd":      "state",
		"appid":    "250",
		"platform": 2,
		"data": map[string]int64{
			"q_recv_source_bytes": 100,
			"q_send_player_bytes": 100,
		},
	}
	state4 = map[string]interface{}{
		"cmd":      "state",
		"appid":    "251",
		"platform": 2,
		"data": map[string]int64{
			"q_recv_source_bytes": 100,
			"q_send_player_bytes": 100,
		},
	}
	//body map[string]interface{} = peerlist
	//list = []map[string]interface{}{login, keepalive, cache,state}
	list = []map[string]interface{}{state3, state4}
)

func Peerclient(ip string, port int) {
	remoteip := net.ParseIP(ip)
	rAddr := &net.UDPAddr{IP: remoteip, Port: port}
	logs.Info("peerclient() start", rAddr)
	conn, err := net.DialUDP("udp", nil, rAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	body, _ := json.Marshal(login)
	if _, err := conn.Write([]byte(body)); err != nil {
		fmt.Println(err)
		return
	}
	msg := make([]byte, 1024)
	n, _ := conn.Read(msg)
	fmt.Println(string(msg[:n]))

	go func() {
		for {
			select {
			case <-time.After(2 * time.Second):
				body, _ := json.Marshal(keepalive)
				if _, err := conn.Write([]byte(body)); err != nil {
					fmt.Println(err)
					return
				}

				msg := make([]byte, 1024)
				n, _ := conn.Read(msg)
				fmt.Println(string(msg[:n]))
			}
		}
	}()

	time.Sleep(5 * time.Second)
	body, _ = json.Marshal(cache1)
	if _, err := conn.Write([]byte(body)); err != nil {
		fmt.Println(err)
		return
	}
	msg = make([]byte, 1024)
	n, _ = conn.Read(msg)
	fmt.Println(string(msg[:n]))
	fmt.Println("----")

	body, _ = json.Marshal(delete_cache)
	if _, err := conn.Write([]byte(body)); err != nil {
		fmt.Println(err)
		return
	}
	msg = make([]byte, 1024)
	n, _ = conn.Read(msg)
	fmt.Println(string(msg[:n]))
	fmt.Println("----")
	time.Sleep(1 * time.Minute)
}

//"183.60.143.82:3030"
func Tcpclient(port int) {
	i := 0
	wg := sync.WaitGroup{}
	for i < 2 {
		wg.Add(1)
		i++
		go func(i int) {
			defer wg.Add(-1)
			conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
			if err != nil {
				fmt.Println(err)
				return
			}
			defer conn.Close()
			msgNum := 0
			for msgNum < 100000 {
				msgNum++
				//body, _ := json.Marshal(keepalive)
				fmt.Printf("%d号udp消息:%d \n", i, msgNum)
				//if _, err := conn.Write([]byte(body)); err != nil {
				//	fmt.Println(err)
				//	return
				//}
				//	time.Sleep(10 * time.Second)
				//msg := make([]byte, 1024)
				//conn.Read(msg)
				for _, msg := range list {
					body, _ := json.Marshal(msg)
					if _, err := conn.Write([]byte(body)); err != nil {
						fmt.Println(err)
						return
					}
					msg := make([]byte, 1024)
					conn.Read(msg)
				}
			}
		}(i)
	}
	wg.Wait()
}

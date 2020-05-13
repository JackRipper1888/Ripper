package dbkit

import (
	"fmt"
	"github.com/Ripper/ctxkit"
	"testing"
)

const (
	USER = "user"
	//[mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	MONGO_URL = "mongodb://yanfa:ff2019yqcsdew@183.60.143.82:36068"
	DBNAME    = "p2p_state"
	MAX_CONN  = 0
)

func TestMongoDB(t *testing.T) {
	//释放心跳协程
	defer ctxkit.CancelAll()
	fmt.Println("TestMongoDB start")
	InitMongoDB(USER, MONGO_URL, DBNAME, MAX_CONN)
	conn := GetMongoConn(USER)

	params := map[string]interface{}{"uid": "asdhasianjkaihs"}
	document := map[string]interface{}{
		//"name":"王大锤",
		"all_stream": []map[string]interface{}{
			{"streamid": "wusha",
				"resourceurl":  "12g3h1ihd",
				"resourcename": "误杀",
			},
		},
		"loc": [2]float32{100.00, 87.55},
	}
	err := conn.Update("p2p", document, params)
	if err != nil {
		fmt.Println(err)
		return
	}

	Geodata, err := conn.FindGeoNear("p2p", 100.00, 13.00, 100, params)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("reading", Geodata)
}

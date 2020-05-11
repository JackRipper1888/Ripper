package utils

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

const (
	USER = "user"
	//[mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	MONGO_URL   = "mongodb://yanfa:ff2019yqcsdew@183.60.143.82:36068"
	DBNAME      = "p2p_state"
	colelection = "p2p"
	MAX_CONN    = 0
)

type UserInfo struct {
	Name      string     `bson:"name"`
	AllStream []string   `bson:"all_stream"`
	Loc       [2]float32 `bson:"loc"`
}

func TestMongoDB(t *testing.T) {
	fmt.Println("TestMongoDB start")
	InitMongoDB(USER, MONGO_URL, DBNAME, MAX_CONN)
	//data := UserInfo{
	//	Name: "王多余",
	//	AllStream:[]string{"wush","wanming"},
	//	Loc: [2]float32{12.33, 44.55},
	//}
	//err := GetMongoConn(USER).Insert(colelection,data)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	Geodata, err := GetMongoConn(USER).FindGeoNear(colelection, 12.00, 13.00, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	if Geodata == nil {
		fmt.Println("reading", Geodata)
		return
	}
	fmt.Println("reading")
	for _, v := range Geodata.([]interface{}) {
		fmt.Println(v.(bson.M)["obj"].(bson.M)["loc"])
	}
}

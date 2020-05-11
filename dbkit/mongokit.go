package utils

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/Ripper/ctxkit"
	"github.com/astaxie/beego/logs"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	SOCKET_TIME_OUT_MS = 10000 //操作超时,默认10s
	CONN_TIME_OUT_MS   = 10000 //连接超时,默认10s
	MAX_POOL_SIZE      = 100   //连接池大小
	MAX_RETRIES        = 5     //连接失败后,重试次数
)

var (
	mongoPool       = make(map[string]*mongoConn)
	mongoPoolRWLock = new(sync.RWMutex)
)

type mongoConn struct {
	globalSession *mgo.Session
	mongoDBName   string
}

func InitMongoDB(mongoConnName, connUrl, dbName string, maxConn int) {
	if connUrl == "" || dbName == "" {
		panic("conn url or db name is empty!")
	}
	if _, isExist := mongoPool[mongoConnName]; isExist {
		panic("mongoConn is exist!")
	}
	fullUrl := setConnUrlOptions(connUrl + "/" + dbName)
	mySession, err := mgo.Dial(fullUrl)
	if err != nil {
		panic(err)
	}
	mySession.SetMode(mgo.Monotonic, true)
	if maxConn == 0 {
		mySession.SetPoolLimit(MAX_POOL_SIZE)
	} else {
		mySession.SetPoolLimit(maxConn)
	}

	monConn := new(mongoConn)
	monConn.mongoDBName = dbName
	monConn.globalSession = mySession
	mongoPoolRWLock.RLock()
	mongoPool[mongoConnName] = monConn
	mongoPoolRWLock.RUnlock()
	//要在主协程中开启 defer ctxkit.CancelAll()
	ctx, _ := ctxkit.CtxAdd()
	go keepAlive(mySession, ctx)
}

func InitMongoDBWithSSL(mongoConnName, connUrl string, dbName string) {
	if connUrl == "" || dbName == "" {
		panic("conn url or db name is empty!")
	}
	if _, isExist := mongoPool[mongoConnName]; isExist {
		panic("mongoConn is exist!")
	}
	fullUrl := setConnUrlOptions(connUrl + "/" + dbName)
	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = true
	dialInfo, err := mgo.ParseURL(fullUrl)
	if err != nil {
		panic(err)
	}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	mySession, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}
	err = mySession.Ping()
	if err != nil {
		panic(err)
	}
	monConn := new(mongoConn)
	monConn.mongoDBName = dbName
	monConn.globalSession = mySession
	mongoPoolRWLock.RLock()
	mongoPool[mongoConnName] = monConn
	mongoPoolRWLock.RUnlock()
	// 定时 Ping，维持长连接
	ctx, _ := ctxkit.CtxAdd()
	go keepAlive(mySession, ctx)
}

func GetMongoConn(mongoConnName string) *mongoConn {
	if mongoConnName == "" {
		panic("mongoConnName is empty!")
	}
	mongoPoolRWLock.RLock()
	defer mongoPoolRWLock.RUnlock()
	v, isExist := mongoPool[mongoConnName]
	if !isExist {
		panic("mongoConn is empty!")
	}
	v.globalSession.Refresh()
	var isSessionOk = true
	err := v.globalSession.Ping()
	if err != nil {
		logs.Error(fmt.Sprintf("globalSession ping fail %v", err))
		isSessionOk = false
		v.globalSession.Refresh()
		for i := 0; i < MAX_RETRIES; i++ {
			err = v.globalSession.Ping()
			if err == nil {
				isSessionOk = true
				logs.Info("Reconnect to mongodb successful.")
				break
			} else {
				logs.Error(fmt.Sprintf("Reconnect to mongodb fail:%v", i))
			}
		}
	}
	if !isSessionOk {
		panic("Reconnect to mongodb fail!")
	}
	return v
}

//插入记录
/*
colelection:表
data:存储文档
*/
func (this *mongoConn) Insert(colelection string, data interface{}) error {
	defer this.globalSession.Close()
	err := this.globalSession.DB(this.mongoDBName).C(colelection).Insert(data)
	this.globalSession.DB(this.mongoDBName).C(colelection).DropIndex()
	if err != nil {
		return err
	}
	return nil
}

//查询地理坐标范围内的数据
/*
colelection:表
pointX:查询地理的横坐标
pointY:查询地理的纵坐标
maxDistance:查询地理范围
*/
func (this *mongoConn) FindGeoNear(colelection string, pointX, pointY, maxDistance float64) (interface{}, error) {
	defer this.globalSession.Close()
	resp := bson.M{}
	db := this.globalSession.DB(this.mongoDBName)
	err := db.Run(bson.D{
		{"geoNear", colelection},
		{"spherical", true},
		{"near", [2]float64{pointX, pointY}},
		//{"num", 3},
		//{"query": { "category": "public" }},
		{"maxDistance", maxDistance},
	}, &resp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return resp["results"], nil
}

func (this *mongoConn) MongoUpdate(colelection string, whereData, updateDate map[string]interface{}) error {
	this.globalSession.DB(this.mongoDBName)
	change := mgo.Change{
		Update:    bson.M{"$set": whereData},
		ReturnNew: false,
		Remove:    false,
		Upsert:    true,
	}
	c := this.globalSession.DB(this.mongoDBName).C(colelection)
	_, err := c.Find(updateDate).Apply(change, nil)
	return err
}
func setConnUrlOptions(connUrl string) string {
	opts := make([]string, 0)
	opts = append(opts, connUrl)
	opts = append(opts, "?")
	opts = append(opts, "authMechanism=SCRAM-SHA-1")
	opts = append(opts, "&maxPoolSize=100")
	return strings.Join(opts, "")
}
func keepAlive(s *mgo.Session, ctx context.Context) {
	c := time.Tick(1 * time.Minute)
	for {
		select {
		case <-c:
			if err := s.Ping(); err != nil {
				logs.Error(err)
			}
		case <-ctx.Done():
			return
		}
	}
}

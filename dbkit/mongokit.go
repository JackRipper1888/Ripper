package dbkit

import (
	"context"
	"crypto/tls"
	"encoding/json"
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

/**
注意 先要将索引要设置成 2dsphere
collection:集合
index:为索引名
db.collection.ensureIndex({index: "2dsphere"})
*/

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
colelection:集合
document:文档
*/
func (this *mongoConn) Insert(colelection string, document interface{}) error {
	conn := this.globalSession.Clone()
	defer conn.Close()
	err := conn.DB(this.mongoDBName).C(colelection).Insert(document)
	if err != nil {
		return err
	}
	return nil
}

//修改记录
/*
colelection:集合
params:筛选条件
document:文档
$:有记录修改,没记录添加
*/
func (this *mongoConn) Update(colelection string, document interface{}, params map[string]interface{}) error {
	conn := this.globalSession.Clone()
	defer conn.Close()
	c := conn.DB(this.mongoDBName).C(colelection)
	change := mgo.Change{
		Update:    bson.M{"$set": document},
		ReturnNew: false,
		Remove:    false,
		Upsert:    true,
	}
	_, err := c.Find(params).Apply(change, nil)
	if err != nil {
		return err
	}
	return nil
}

//查询地理坐标范围内的数据
/*
colelection:集合
pointX:查询地理的横坐标
pointY:查询地理的纵坐标
maxDistance:查询地理范围
*/
func (this *mongoConn) FindGeoNear(colelection string, pointX, pointY, maxDistance float64, queryParams map[string]interface{}) ([]map[string]interface{}, error) {
	conn := this.globalSession.Clone()
	defer conn.Close()
	resp := bson.M{}
	db := conn.DB(this.mongoDBName)
	var reqList = []bson.DocElem{
		{"geoNear", colelection},
		//是否用球面来计算距离，如果是2dsphere必须为true
		{"spherical", true},
		//指定附近点的坐标 对于2dsphere用GeoJson，对于2s用坐标对
		{"near", [2]float64{pointX, pointY}},
		//限制的最大距离，如果是GeomJson单位为米，如果是坐标对单位为弧度
		{"maxDistance", maxDistance},
		//对返回的基于距离的结果，乘以这个算子
		{"distanceMultiplier", 6371},
		//限制的最小距离，如果是GeomJson单位为米，如果是坐标对单位为弧度
		//{"minDistance", 3},
		//返回的最大数，默认是100
		//{"limit", 3},
		//{"num", 3},
	}
	for k, v := range queryParams {
		reqList = append(reqList, bson.DocElem{"query", bson.M{k: v}})
	}

	err := db.Run(reqList, &resp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	body, err := bson.MarshalJSON(resp["results"])
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var jsonDocuments []map[string]interface{}
	if err = json.Unmarshal(body, &jsonDocuments); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return jsonDocuments, nil
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

package utils

import (
	"fmt"
	"github.com/Ripper/errkit"
	"github.com/go-redis/redis"
	"strings"
	"sync"
	"time"
)

var (
	redisList                     = make(map[string]*RedisClient)
	redisListRWLock *sync.RWMutex = new(sync.RWMutex)
)

type RedisClient struct {
	client      *redis.Client
	redisInited bool
}

func GetRedisClient(rc string) *RedisClient {
	redisListRWLock.RLock()
	defer redisListRWLock.RUnlock()
	if !redisList[rc].redisInited || redisList[rc].client == nil {
		return nil
	}
	return redisList[rc]
}

func GetRedisPipe(rc string) *redis.Client {
	redisListRWLock.RLock()
	defer redisListRWLock.RUnlock()
	if !redisList[rc].redisInited || redisList[rc].client == nil {
		return nil
	}
	return redisList[rc].client
}

//初始化数据库
/*
@addr 地址
@passwd 密码
@dbNum 连接号
@maxConn 最大连接数
*/
func InitRedis(RedisClientName string, addr, passwd string, dbNum, maxConn int) {
	redisListRWLock.Lock()
	defer redisListRWLock.Unlock()
	if redisList[RedisClientName] != nil {
		return
	}
	if addr == "" {
		panic("redis addr is empty!")
	}
	if dbNum < 0 || dbNum > 16 {
		panic("redis dbNum is error!")
	}
	if maxConn == 0 {
		maxConn = 10
	}
	rc := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     passwd,           //如果没有密码,默认为空
		DB:           dbNum,            //默认选择0数据库
		MaxRetries:   3,                //连接失败后重试3次
		DialTimeout:  10 * time.Second, //拨号超时
		WriteTimeout: 5 * time.Second,  //写超时
		PoolSize:     maxConn,          //最大连接数
		IdleTimeout:  200 * time.Second,
	})
	pong, err := rc.Ping().Result()
	if err != nil {
		panic(fmt.Sprintf("redis ping is error %v", err))
	}
	if strings.ToUpper(pong) != "PONG" {
		panic("redis conn return is not pong")
	}
	var data = RedisClient{
		client:      rc,
		redisInited: true,
	}
	redisList[RedisClientName] = &data
}

// 向数据库中添加键值对内容
/*
@key 	主键
@value 	内容
@sec	过期时间,单位秒,0:永不过期
*/
func (rc *RedisClient) RedisSetWithExpire(key string, value string, sec time.Duration) error {
	if key == "" || value == "" {
		return errkit.New(-1, "redis params is empty")
	}
	err := rc.client.Set(key, value, sec).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rc *RedisClient) RedisGet(key string) (string, error) {
	if key == "" {
		return "", errkit.New(-1, "redis params is empty")
	}
	Result, err := rc.client.Get(key).Result()
	if err == redis.Nil {
		return "", errkit.New(100018, "redis get key does not exists")
	} else if err != nil {
		return "", err
	}
	return Result, nil
}

// 向数据库中添加键值对内容,值是一组set集合
/*
@key 	主键
@fields 	内容
@sec		过期时间,单位秒,0:永不过期
*/
func (rc *RedisClient) RedisSetSAddWithExpire(key string, sec time.Duration, fields []interface{}) error {
	if key == "" || fields == nil {
		return errkit.New(-1, "redis params is empty")
	}
	err := rc.client.SAdd(key, fields...).Err()
	if err != nil {
		return err
	}
	if sec > 0 { //设置KEY过期时间
		rc.client.Expire(key, sec)
	}
	return nil
}

// 从数据库中移除并返回集合中的多个随机元素
// @key 	主键
// @count 	数量
func (rc *RedisClient) RedisSPopN(key string, count int64) ([]string, error) {
	if key == "" {
		return nil, errkit.New(-1, "redis params is empty")
	}
	result, err := rc.client.SPopN(key, count).Result()
	if err == redis.Nil {
		return nil, errkit.New(100018, "redis get key does not exists")
	} else if err != nil {
		return nil, err
	}
	return result, nil
}

// 从数据库中返回集合中的多个随机元素
// @key 	主键
// @count 	数量
func (rc *RedisClient) SRandMemberN(key string, count int64) ([]string, error) {
	if key == "" {
		return nil, errkit.New(-1, "redis params is empty")
	}
	result, err := rc.client.SRandMemberN(key, count).Result()
	if err == redis.Nil {
		return nil, errkit.New(100018, "redis get key does not exists")
	} else if err != nil {
		return nil, err
	}
	return result, nil
}

// 从数据库中返回集合
// @key 	主键
// @count 	数量
func (rc *RedisClient) SMembers(key string) ([]string, error) {
	if key == "" {
		return nil, errkit.New(-1, "redis params is empty")
	}
	result, err := rc.client.SMembers(key).Result()
	if err == redis.Nil {
		return nil, errkit.New(100018, "redis get key does not exists")
	} else if err != nil {
		return nil, err
	}
	return result, nil
}

func (rc *RedisClient) HMSet(key string, data map[string]interface{}, sec time.Duration) error {
	if key == "" {
		return errkit.New(-1, "redis params is empty")
	}
	err := rc.client.HMSet(key, data).Err()
	if err != nil {
		return err
	}
	if sec > 0 { //设置KEY过期时间
		rc.client.Expire(key, sec)
	}
	return nil
}

func (rc *RedisClient) HMGetAll(key string) (map[string]string, error) {
	if key == "" {
		return nil, errkit.New(-1, "redis params is empty")
	}
	MapData, err := rc.client.HGetAll(key).Result()
	if err == redis.Nil {
		return nil, errkit.New(100018, "redis get key does not exists")
	} else if err != nil {
		return nil, err
	}

	return MapData, nil
}

// 从数据库获取对应键内容
// @key 	主键
func (rc *RedisClient) HMGetVal(key string, value ...string) ([]interface{}, error) {
	if key == "" {
		return nil, errkit.New(-1, "redis params is empty")
	}
	result, err := rc.client.HMGet(key, value...).Result()
	if err == redis.Nil {
		return nil, errkit.New(100018, "redis get key does not exists")
	} else if err != nil {
		return result, err
	} else {
		return result, nil
	}
}

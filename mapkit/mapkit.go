package mapkit

import (
	"fmt"
	"sync"
)

// 总的map
type ConcurrentMap []*ConcurrentMapShared

// 默认分片数
var (
	SHARE_COUNT int = 64
)

// 单个map分片
type ConcurrentMapShared struct {
	items map[string]interface{} // 本分片内的map
	mu    sync.RWMutex           // 本分片的专用锁
}

// 新建一个map
func NewConcurrentMap(share_cont int) *ConcurrentMap {
	SHARE_COUNT = share_cont
	m := make(ConcurrentMap, SHARE_COUNT)
	for i := 0; i < SHARE_COUNT; i++ {
		m[i] = &ConcurrentMapShared{
			items: map[string]interface{}{},
			mu:    sync.RWMutex{},
		}
	}
	return &m
}

// GetSharedMap 获取key对应的map分片
func (m ConcurrentMap) GetSharedMap(key string) *ConcurrentMapShared {
	return m[uint(fnv32(key))%uint(SHARE_COUNT)]
}

// hash函数
func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	prime32 := uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

// Set 设置key,value
func (m ConcurrentMap) Set(key string, value interface{}) {
	sharedMap := m.GetSharedMap(key) // 找到对应的分片map
	sharedMap.mu.Lock()              // 加锁(全锁定)
	sharedMap.items[key] = value     // 赋值
	sharedMap.mu.Unlock()            // 解锁
}

// Get 获取key对应的value
func (m ConcurrentMap) Get(key string) (value interface{}, ok bool) {
	sharedMap := m.GetSharedMap(key) // 找到对应的分片map
	fmt.Println(uint(fnv32(key)) % uint(SHARE_COUNT))
	sharedMap.mu.RLock()             // 加锁(读锁定)
	value, ok = sharedMap.items[key] // 取值
	sharedMap.mu.RUnlock()           // 解锁
	return value, ok
}

// Get 获取key对应的value删除对应的data
func (m ConcurrentMap) Delete(key string) {
	sharedMap := m.GetSharedMap(key) // 找到对应的分片map
	sharedMap.mu.Lock()
	delete(sharedMap.items, key)
	sharedMap.mu.Unlock()
}

// Get 获取key对应的value删除对应的data
func (m ConcurrentMap) DeleteInBatches(keys ...string) {
	keylist := make(map[uint][]string, 0)
	for _, key := range keys {
		k := uint(fnv32(key)) % uint(SHARE_COUNT)
		keylist[k] = append(keylist[k], key)
	}
	for i, mkey := range keylist {
		sharedMap := m[i] // 找到对应的分片map
		sharedMap.mu.Lock()
		for _, key := range mkey {
			delete(sharedMap.items, key)
		}
		sharedMap.mu.Unlock()
	}
}

// Count 统计key个数
func (m ConcurrentMap) Count() int {
	count := 0
	for i := 0; i < SHARE_COUNT; i++ {
		m[i].mu.RLock() // 加锁(读锁定)
		count += len(m[i].items)
		m[i].mu.RUnlock() // 解锁
	}
	return count
}

// Keys 所有的key方法(方法:开多个协程分别对分片map做统计再汇总 优点:量大时,阻塞时间较短)
func (m ConcurrentMap) Keys() []string {
	count := m.Count()
	keys := make([]string, count)

	ch := make(chan string, count) // 通道,遍历时
	// 单独起一个协程
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(SHARE_COUNT)

		for i := 0; i < SHARE_COUNT; i++ {
			// 每个分片map,单独起一个协程进行统计
			go func(ms *ConcurrentMapShared) {
				defer wg.Done()

				ms.mu.RLock() // 加锁(读锁定)
				for k := range ms.items {
					ch <- k // 压入通道
				}
				ms.mu.RUnlock() // 解锁
			}(m[i])
		}

		// 等待所有协程执行完毕
		wg.Wait()
		close(ch) // 一定要关闭通道,因为不关闭的话,后面的range不会结束!!!
	}()

	// 遍历通道,压入所有的key
	for k := range ch {
		keys = append(keys, k)
	}
	return keys
}

func (m ConcurrentMap) Range(f func(k string, v interface{}) bool) {
	for _, mp := range m {
		mp.mu.RLock()
		for k, v := range mp.items {
			if !f(k, v) {
				return
			}
		}
		mp.mu.RUnlock()
	}
}

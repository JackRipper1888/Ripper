package bufkit

import (
	"container/list"
	"math/rand"
	"net"
	"time"
)

var (
	free, makes int
)

func makeBuffer() []byte {
	makes += 1
	return make([]byte, rand.Intn(5000000)+5000000)
}

type queued struct {
	when  time.Time
	slice []byte
}

// peer指令
type PeerCmd struct {
	peerAddr net.UDPAddr
	data     []byte
}

func MakeRecycler() (get, give chan []byte) {
	get = make(chan []byte)
	give = make(chan []byte)

	go func() {
		q := new(list.List)
		for {
			if q.Len() == 0 {
				q.PushFront(queued{when: time.Now(), slice: makeBuffer()})
			}

			e := q.Front()
			timeout := time.NewTimer(time.Minute)
			select {
			case b := <-give:
				timeout.Stop()
				q.PushFront(queued{when: time.Now(), slice: b})

			case get <- e.Value.(queued).slice:
				timeout.Stop()
				q.Remove(e)

			case <-timeout.C:
				e := q.Front()
				for e != nil {
					n := e.Next()
					if time.Since(e.Value.(queued).when) > time.Minute {
						q.Remove(e)
						e.Value = nil
					}
					e = n
				}
			}
		}
	}()
	return
}

package pubsub

import (
	"sync"
	"time"
)

type (
	Subscriber chan interface{}
	TopicFunc  func(v interface{}) bool
)
type Publisher struct {
	subscribers map[Subscriber]TopicFunc
	buffersize  int
	timeout     time.Duration
	sync.RWMutex
}

func NewPublisher(publishTimeout time.Duration, bufferSize int) *Publisher {
	return &Publisher{
		buffersize:  bufferSize,
		timeout:     publishTimeout,
		subscribers: make(map[Subscriber]TopicFunc),
	}
}

func (p *Publisher) Publish(v interface{})  {
	p.RLock()
	defer p.RUnlock()
	var wg sync.WaitGroup
	for sub,topic := range p.subscribers{
		wg.Add(1)
		go p.sendTopic(sub,topic,v,wg)
	}
	wg.Wait()
}

func (p *Publisher)sendTopic(sub Subscriber,topic TopicFunc,v interface{},wg sync.WaitGroup)  {
	defer wg.Done()
	if topic != nil && !topic(v){
		return
	}
	select{
	case sub <- v:
	case <- time.After(p.timeout):
	}
}
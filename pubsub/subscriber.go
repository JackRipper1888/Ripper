package pubsub

func (p *Publisher) SubscribeTopic(topic TopicFunc) Subscriber {
	ch := make(Subscriber, p.buffersize)
	p.Lock()
	p.subscribers[ch] = topic
	p.Unlock()
	return ch
}
func (p *Publisher)Evict(sub Subscriber) {
	p.Lock()
	defer p.Unlock()
	delete(p.subscribers, sub)
	close(sub)
}
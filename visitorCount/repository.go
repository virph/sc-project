package visitorCount

type RedisVisitorCountRepository interface {
	Increase()
	Get() int
}

type NsqVisitorCountRepository interface {
	Publish(topic string, payload []byte) error
	AddVisitor()
}

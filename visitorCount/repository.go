package visitorCount

type RedisVisitorCountRepository interface {
	Increase()
	Get() int
}

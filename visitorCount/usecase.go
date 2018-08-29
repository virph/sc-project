package visitorCount

type VisitorCountUsecase interface {
	Get() int
	Increase()
	PublishIncrease()
}

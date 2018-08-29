package usecase

import (
	"log"

	"github.com/virph/sc-project/visitorCount"
)

type visitorCountUsecase struct {
	repository visitorCount.RedisVisitorCountRepository
	nsq        visitorCount.NsqVisitorCountRepository
}

func NewVisitorCountUsecase(repository *visitorCount.RedisVisitorCountRepository, nsq *visitorCount.NsqVisitorCountRepository) visitorCount.VisitorCountUsecase {
	usecase := visitorCountUsecase{
		repository: *repository,
		nsq:        *nsq,
	}
	return &usecase
}

func (u *visitorCountUsecase) Get() int {
	return u.repository.Get()
}

func (u *visitorCountUsecase) PublishIncrease() {
	log.Println("usecase increase called")
	u.nsq.AddVisitor()
}

func (u *visitorCountUsecase) Increase() {
	u.repository.Increase()
}

package usecase

import (
	"github.com/virph/sc-project/models"
	"github.com/virph/sc-project/user"
)

type userUsecase struct {
	repository user.PostgreUserRepository
}

func NewUserUsecase(repository *user.PostgreUserRepository) user.UserUsecase {
	usecase := userUsecase{
		repository: *repository,
	}

	return &usecase
}

func (u *userUsecase) Find(name string) []models.User {
	return u.repository.Find(name)
}

package user

import "github.com/virph/sc-project/models"

type UserUsecase interface {
	Find(name string) []models.User
}

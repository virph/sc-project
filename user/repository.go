package user

import "github.com/virph/sc-project/models"

type PostgreUserRepository interface {
	Find(name string) []models.User
}

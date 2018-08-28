package models

import (
	"time"

	"github.com/lib/pq"
)

type User struct {
	ID         string      `db:"user_id"`
	Name       string      `db:"full_name"`
	MSISDN     string      `db:"msisdn"`
	Email      string      `db:"user_email"`
	BirthDate  pq.NullTime `db:"birth_date"`
	CreateTime pq.NullTime `db:"create_time"`
	UpdateTime pq.NullTime `db:"update_time"`
}

func (u *User) Age() int {
	if !u.BirthDate.Valid {
		return 0
	}
	return int(time.Since(u.BirthDate.Time).Hours() / 24 / 365)
}

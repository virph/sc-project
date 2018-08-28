package repository

import (
	"database/sql"
	"log"

	"github.com/virph/sc-project/models"
	"github.com/virph/sc-project/user"
)

type postgreUserRepository struct {
	Connection *sql.DB
}

func NewPostgreUserRepository(c *sql.DB) user.PostgreUserRepository {
	return &postgreUserRepository{
		Connection: c,
	}
}

func (r *postgreUserRepository) Find(name string) []models.User {
	users := make([]models.User, 0)
	st, err := r.Connection.Prepare(`
		select
			coalesce(user_id, 0),
			coalesce(full_name, '-'),
			coalesce(msisdn, '-'),
			coalesce(user_email, '-'),
			birth_date,
			create_time,
			update_time
		from
			public.ws_user
		where
			full_name like $1
		limit 10`)
	if err != nil {
		log.Println(err)
	}

	rows, err := st.Query("%" + name + "%")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		u := models.User{}
		err = rows.Scan(
			&u.ID,
			&u.Name,
			&u.MSISDN,
			&u.Email,
			&u.BirthDate,
			&u.CreateTime,
			&u.UpdateTime,
		)
		if err != nil {
			log.Println(err)
		}
		users = append(users, u)
	}
	return users
}

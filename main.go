package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	redigo "github.com/garyburd/redigo/redis"
	_ "github.com/lib/pq"
	"github.com/virph/sc-project/user/delivery"
	userrepo "github.com/virph/sc-project/user/repository"
	userusecase "github.com/virph/sc-project/user/usecase"
	"github.com/virph/sc-project/visitorCount/repository"
)

const (
	host     = "devel-postgre.tkpd"
	port     = 5432
	user     = "tkpdtraining"
	password = "1tZCjrIcYeR1uQblQz0gBlIFU"
	dbname   = "tokopedia-dev-db"

	serverPort = ":3000"
)

func initDb() *sql.DB {
	var (
		db  *sql.DB
		err error
	)

	cs := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, password, host, port, dbname)
	db, err = sql.Open("postgres", cs)
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Database connected")
	return db
}

func initRedis() *redigo.Pool {
	return redigo.NewPool(func() (r redigo.Conn, err error) {
		return redigo.Dial("tcp", "devel-redis.tkpd:6379")
	}, 1)
}

func main() {
	db := initDb()
	defer db.Close()

	redisPool := initRedis()
	defer redisPool.Close()

	redisVisitorCountRepository := repository.NewRedisVisitorCountRepository(redisPool)

	redisVisitorCountRepository.Increase()
	fmt.Println("redis", redisVisitorCountRepository.Get())
	redisVisitorCountRepository.Increase()
	fmt.Println("redis", redisVisitorCountRepository.Get())
	redisVisitorCountRepository.Increase()
	fmt.Println("redis", redisVisitorCountRepository.Get())
	redisVisitorCountRepository.Increase()
	fmt.Println("redis", redisVisitorCountRepository.Get())

	pgUserRepo := userrepo.NewPostgreUserRepository(db)
	userUsecase := userusecase.NewUserUsecase(&pgUserRepo)
	delivery.NewUserHandler(&userUsecase)

	log.Println("Listening on port", serverPort)
	http.ListenAndServe(serverPort, nil)
}

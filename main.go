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
)

const (
	dbHost = "devel-postgre.tkpd"
	dbPort = 5432
	dbUser = "tkpdtraining"
	dbPass = "1tZCjrIcYeR1uQblQz0gBlIFU"
	dbName = "tokopedia-dev-db"

	redisHost = "devel-redis.tkpd:6379"

	serverPort = ":3000"
)

func initDb() *sql.DB {
	var (
		db  *sql.DB
		err error
	)

	cs := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)
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
	pool := redigo.NewPool(func() (r redigo.Conn, err error) {
		return redigo.Dial("tcp", redisHost)
	}, 1)

	con := pool.Get()
	defer con.Close()

	_, err := con.Do("PING")
	if err != nil {
		log.Fatalln("Can't connect to the Redis database")
	}

	log.Println("Redis connected")
	return pool
}

func main() {
	db := initDb()
	defer db.Close()

	redisPool := initRedis()
	defer redisPool.Close()

	pgUserRepo := userrepo.NewPostgreUserRepository(db)
	userUsecase := userusecase.NewUserUsecase(&pgUserRepo)
	delivery.NewUserHandler(&userUsecase)

	log.Println("Listening to port", serverPort)
	http.ListenAndServe(serverPort, nil)
}

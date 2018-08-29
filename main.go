package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/virph/sc-project/visitorCount"

	redigo "github.com/garyburd/redigo/redis"
	_ "github.com/lib/pq"
	nsq "github.com/nsqio/go-nsq"
	"github.com/virph/sc-project/user/delivery"
	userrepo "github.com/virph/sc-project/user/repository"
	userusecase "github.com/virph/sc-project/user/usecase"
	visitorcountrepo "github.com/virph/sc-project/visitorCount/repository"
	visitorcountusecase "github.com/virph/sc-project/visitorCount/usecase"
)

const (
	dbHost = "devel-postgre.tkpd"
	dbPort = 5432
	dbUser = "tkpdtraining"
	dbPass = "1tZCjrIcYeR1uQblQz0gBlIFU"
	dbName = "tokopedia-dev-db"

	redisHost = "devel-redis.tkpd:6379"

	serverPort = ":3000"

	nsqHost = "127.0.0.1:4150"
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

func initNSQ() *nsq.Producer {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(nsqHost, config)
	if err != nil {
		log.Fatalln(err)
	}
	if err := producer.Ping(); err != nil {
		log.Fatalln(err)
	}
	log.Println("NSQ connected")
	return producer
}

func initNSQConsumer(visitorCountUsecase visitorCount.VisitorCountUsecase) {
	config := nsq.NewConfig()
	c, err := nsq.NewConsumer("visitor_count", "visitor_count_channel", config)
	if err != nil {
		log.Fatalln(nil)
	}

	c.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		visitorCountUsecase.Increase()
		return nil
	}))

	err = c.ConnectToNSQD(nsqHost)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	db := initDb()
	defer db.Close()

	redisPool := initRedis()
	defer redisPool.Close()

	nsqProducer := initNSQ()

	pgUserRepo := userrepo.NewPostgreUserRepository(db)
	userUsecase := userusecase.NewUserUsecase(&pgUserRepo)

	nsqVisitorCountRepo := visitorcountrepo.NewNsqVisitorCountRepository(nsqProducer)

	rsVisitorCountRepo := visitorcountrepo.NewRedisVisitorCountRepository(redisPool)
	visitorCountUsecase := visitorcountusecase.NewVisitorCountUsecase(&rsVisitorCountRepo, &nsqVisitorCountRepo)

	delivery.NewUserHandler(&userUsecase, &visitorCountUsecase)

	initNSQConsumer(visitorCountUsecase)

	log.Println("Listening to port", serverPort)
	http.ListenAndServe(serverPort, nil)
}

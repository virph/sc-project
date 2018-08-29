package repository

import (
	redigo "github.com/garyburd/redigo/redis"
	"github.com/virph/sc-project/visitorCount"
)

const (
	key    = "training-sc:wkwkw"
	expire = 30 * 60
)

type redisVisitorCountRepository struct {
	pool *redigo.Pool
}

func NewRedisVisitorCountRepository(redisPool *redigo.Pool) visitorCount.RedisVisitorCountRepository {
	return &redisVisitorCountRepository{
		pool: redisPool,
	}
}

func (r *redisVisitorCountRepository) Increase() {
	con := r.pool.Get()
	defer con.Close()

	con.Do("INCR", key)
	updateExpire(&con)
}

func (r *redisVisitorCountRepository) Get() int {
	con := r.pool.Get()
	defer con.Close()

	v, _ := redigo.Int(con.Do("GET", key))
	return v
}

func updateExpire(con *redigo.Conn) {
	redigo.Conn(*con).Do("EXPIRE", key, expire)
}

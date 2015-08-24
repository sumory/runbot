package common

import (
	"time"
	"github.com/garyburd/redigo/redis"
)

type User struct {
	Uid int32
}

type Relation struct {
	Uid    int32
	Target int32
	Time   float64
}

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func NewRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle: 10,
		MaxActive: 500,
		IdleTimeout: 240 * time.Second,
		Dial: func () (redis.Conn, error) {
			c, err := redis.Dial("tcp",  "192.168.100.185:6379")
			if err != nil {
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Wait:true,
	}
}

func InitRedis() redis.Conn {
	conn, err := redis.Dial("tcp", "192.168.100.185:6379")
	if err != nil {
		panic(err.Error())
	}
	return conn
}

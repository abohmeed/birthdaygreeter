package main

import (
	"log"
	"os"

	"github.com/gomodule/redigo/redis"
)

var host = "127.0.0.1"
var port = "6379"
var password = ""

func main() {
	host = os.Getenv("REDIS_MASTER_HOST")
	port = os.Getenv("REDIS_PORT")
	password = os.Getenv("REDIS_PASSWORD")
	pool := newPool()
	conn := pool.Get()
	defer conn.Close()
	if err := ping(conn); err != nil {
		log.Println(err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
func newPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host+":"+port)
			if err != nil {
				log.Println("Could not reach Redis", err)
			}
			_, err = c.Do("AUTH", password)
			if err != nil {
				log.Println("Could not authenticate to Redis", err)
			}
			return c, err
		},
	}
}

func ping(c redis.Conn) error {
	_, err := redis.String(c.Do("PING"))
	if err != nil {
		return err
	}
	return nil
}

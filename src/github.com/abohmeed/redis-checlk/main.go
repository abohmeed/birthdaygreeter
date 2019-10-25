package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gomodule/redigo/redis"
)

func main() {
	hostPtr := flag.String("h", "127.0.0.1", "[Required] The `Redis Host`")
	portPtr := flag.String("p", "6379", "[Required] The `port` ")
	flag.Parse()
	flag.Usage = func() {
		fmt.Printf("Usage: redis-check -h host  -p port")
		flag.PrintDefaults()
	}
	if *hostPtr == "" || *portPtr == "" {
		flag.Usage()
		os.Exit(1)
	}
	pool := newPool(*hostPtr, *portPtr)
	conn := pool.Get()
	defer conn.Close()
	if err := ping(conn); err != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
func newPool(host string, port string) *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host+":"+port)
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

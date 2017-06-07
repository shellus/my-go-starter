package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

var r *redis.Client

func ExampleNewClient() {
	r = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := r.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
}

func ExampleClient() {
	err := r.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := r.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := r.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exists")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exists
}

func main() {
	ExampleNewClient()
	ExampleClient()
}

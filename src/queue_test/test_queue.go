package main

import (
	"fmt"
	"queue"
)

type MyJob struct {
	queue.Job
	payload string
}
func main() {
	q := queue.NewQueue()

	q.Sub(func(j MyJob) {
		fmt.Println(j.payload)
	})

	q.Pub(MyJob{payload:"hahaha"})
	q.Pub(MyJob{payload:"hahaha2"})
	q.Pub(MyJob{payload:"hahaha3"})
	q.Pub(MyJob{payload:"hahaha4"})
}

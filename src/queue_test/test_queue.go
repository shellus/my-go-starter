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

	q.Push(MyJob{payload:"hahaha"})
	q.Push(MyJob{payload:"hahaha2"})
	q.Push(MyJob{payload:"hahaha3"})
	q.Push(MyJob{payload:"hahaha4"})
}

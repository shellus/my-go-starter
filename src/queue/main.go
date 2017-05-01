package main

import (
	"time"
	"container/list"
)

type Queue struct {
	jobs list.List
}
type Job struct {

}

func NewQueue() Queue {

	return Queue{
		jobs: list.New(),
	}
}
func (q Queue) push(j Job){
	append(q.jobs, j)
}

func (q Queue) sub(f func(j Job)){



	for{
		q.jobs.Front()
		time.Sleep(1 * time.Second)
	}
}


func main() {

}

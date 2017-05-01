package queue

import (
)

type Queue struct {
	jobChan chan Job
}
type Job struct {

}

func NewQueue() Queue {
	return Queue{
		jobChan: make(chan Job),
	}
}

func (q Queue) Push(j Job) {
	q.jobChan <- j
}

func (q Queue) Sub(f func(j Job)) {
	for j := range q.jobChan{
		f(j)
	}
}
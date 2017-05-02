package queue

type Queue struct {
	jobs chan Job
	concurrent chan bool
	subscriber func(j Job)
}
type Job struct {
	Value interface{}
}

func NewQueue(concurrentNumber int) (q Queue) {
	q = Queue{
		concurrent: make(chan bool, concurrentNumber),
		jobs: make(chan Job, 10000),
	}
	return
}

func (q *Queue) Push(j Job){
	q.jobs <- j
}

func (q *Queue) Sub(f func(j Job)){
	q.subscriber = f
}
func (q *Queue) Work(){
	for j := range q.jobs{
		q.concurrent <- true
		go q.call(j)
	}
}
func (q *Queue) call (j Job) {
	q.subscriber(j)
	<- q.concurrent
}

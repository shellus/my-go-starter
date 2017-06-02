package queue

// go的任务队列，或者说goroutine管理器
// 创建一个queue(NewQueue)并设置并发数，绑定一个消费者(Sub)，然后Push一堆任务。然后Work阻塞执行。执行完了之后，就会退出Work方法。

type queue struct {
	jobs       chan Job
	concurrent chan bool
	subscriber func(j Job)
}

type Job struct {
	Value interface{}
}

func NewQueue(concurrentNumber int, channel string) (q *queue) {
	// channel is back
	q = &queue{
		concurrent: make(chan bool, concurrentNumber),
		jobs: make(chan Job, 10000),
	}
	return
}

func (q *queue) Push(j Job) {
	q.jobs <- j
}

func (q *queue) Sub(f func(j Job)) {
	q.subscriber = f
}
func (q *queue) Work() {
	L:
	for {
		select {
		case j := <-q.jobs:
			q.concurrent <- true
			go q.call(j)
		default:
			for i := 0; i < cap(q.concurrent); i++ {
				q.concurrent <- true
			}
			break L
		}
	}

}
func (q *queue) call(j Job) {
	defer func() {
		<-q.concurrent
	}()
	q.subscriber(j)
}

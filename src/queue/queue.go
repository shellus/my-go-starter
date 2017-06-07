package queue

import (
	"fmt"
	"github.com/go-redis/redis"
	"encoding/json"
	"time"
)

const prefix string = "queue"
// go的任务队列，或者说goroutine管理器
// 创建一个queue(NewQueue)并设置并发数，绑定一个消费者(Sub)，然后Push一堆任务。然后Work阻塞执行。执行完了之后，就会退出Work方法。

type queue struct {
	channelName string
	concurrent  chan bool
	subscriber  func(j Job) (err error)
	redis       *redis.Client
}

type Job struct {
	Value interface{}
}

func NewQueue(concurrentNumber int, channelName string) (q *queue) {
	// channel is back
	q = &queue{
		channelName: channelName,
		concurrent: make(chan bool, concurrentNumber),
		redis: redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "", // no password set
			DB:       0, // use default DB
		}),
	}
	return
}

func (q *queue) Push(j *Job) {
	q.redis.LPush(prefix + ":" + q.channelName, serialization(j))
}

func (q *queue) Sub(f func(j Job) (err error)) {
	q.subscriber = f
}
func (q *queue) Work() {
	for {
		s, err := q.redis.BRPopLPush(prefix + ":list:" + q.channelName, prefix + ":run:" + q.channelName, time.Minute).Result()
		if err != nil {

		}

		j := &Job{}
		deserialization(s[1], j)
		q.concurrent <- true
		go q.call(j)
	}

}
func (q *queue) call(j *Job) {
	defer func() {
		<-q.concurrent
		// 还是不要拦截这些致命错误了
		//if r := recover(); r != nil {
		//	q.Push(j)
		//	fmt.Println(r)
		//	fmt.Printf("%T", r)
		//}
	}()

	// call
	err := q.subscriber(*j)

	// 如果出错，把job放回去
	if err != nil {
		q.Push(j)
		fmt.Printf("%+v\n", err)
	}

}

func serialization(j *Job) (s string) {
	buf, err := json.Marshal(j)
	if err != nil {
		return ""
	}
	return string(buf)
}

func deserialization(s string, j *Job) (err error) {
	err = json.Unmarshal([]byte(s), j)
	if err != nil {
		return err
	}
	return nil
}
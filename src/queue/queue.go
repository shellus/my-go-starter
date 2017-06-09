package queue

import (
	"fmt"
	"github.com/go-redis/redis"
	"encoding/json"
	"time"
	"github.com/pkg/errors"
)

const prefix string = "queue"

// is old
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

func (q *queue) listKey() (string) {
	return prefix + ":list:" + q.channelName
}
func (q *queue) runKey() (string) {
	return prefix + ":run:" + q.channelName
}

func (q *queue) Pub(j *Job) {
	q.redis.LPush(q.listKey(), serialization(j))
}

func (q *queue) Sub(f func(j Job) (err error)) {
	q.subscriber = f
}

func (q *queue) Work() {
	for {
		s, err := q.redis.BRPopLPush(q.listKey(), q.runKey(), time.Minute).Result()
		if err != nil {
			if err == errors.New("redis: nil") {
				continue
			}
			panic(err)
		}

		j := &Job{}

		if err := deserialization(s, j); err != nil {
			fmt.Println(errors.New(err.Error()))

		}
		q.concurrent <- true
		go q.call(j)
	}

}
/**
运行中的，全部撤回待运行列表。
 */
//func Restart()(r int) {
//
//}
func (q *queue) rollback(j *Job) {
	pipe := q.redis.TxPipeline()

	// 从运行中列表移除job
	pipe.LRem(q.runKey(), 1, serialization(j))

	// 如果出错，把job放回去
	pipe.LPush(q.listKey(), serialization(j))

	_, _ = pipe.Exec()
}
func (q *queue) done(j *Job) {
	// 从运行中列表移除job
	q.redis.LRem(q.runKey(), 1, serialization(j))

}
func (q *queue) call(j *Job) {
	defer func() {
		<-q.concurrent
	}()

	// call
	err := q.subscriber(*j)

	if err != nil {
		fmt.Printf("%+v\n", err)
		q.rollback(j)
	} else {
		q.done(j)
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
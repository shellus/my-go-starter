package rand_reader

import (
	"io"
	"math/rand"
	"time"
)

// 模拟网络IO的reader，哈哈，每次随机返回一点read。

type Reader struct {
	s []byte
	i int // current reading index
}
func init(){
	rand.Seed(time.Now().Unix())
}

func random(min, max int) int {
	return rand.Intn(max - min) + min
}

func New(b []byte) *Reader {
	return &Reader{b, 0}
}
func (r *Reader) Write(b []byte){
	r.s = append(r.s, b[:]...)
}

func (r *Reader) Read(b []byte) (n int, err error) {
	if r.i >= len(r.s) {
		return 0, io.EOF
	}

	rd := random(20, 1024)

	// 如果这次随机超出，就设置为取完
	if r.i + rd > len(r.s) {
		rd = len(r.s) - r.i
	}

	n = copy(b, r.s[r.i:r.i + rd])
	r.i += n
	return
}
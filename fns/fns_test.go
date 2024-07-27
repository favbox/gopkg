package fns

import (
	"fmt"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bytedance/gopkg/lang/channel"
	"github.com/stretchr/testify/assert"
)

func TestToMap(t *testing.T) {
	ss := []string{"a", "b", "c"}
	r := ToMap[int, string](ss, func(s string) int {
		return slices.Index(ss, s)
	})
	fmt.Println(r)
}

func TestFilterParam(t *testing.T) {
	ss := []string{"a", "b", "c", "bc"}
	rtn := FilterParam[string, string](ss, "b", func(s string, s2 string) bool {
		return strings.HasPrefix(s, s2)
	})
	fmt.Println(rtn)
}

func TestSelectRandom(t *testing.T) {
	ss := []string{"a", "b", "c", "bc"}
	for i := 0; i < 5; i++ {
		fmt.Println(SelectRandom(ss, 2))
	}
}

func BenchmarkSelectRandom(b *testing.B) {
	b.ReportAllocs()
	ss := []string{"a", "b", "c", "bc"}
	for i := 0; i < b.N; i++ {
		SelectRandom(ss, 2)
	}
}

func TestUnique(t *testing.T) {
	ss := []string{"a", "b", "c", "b", "a", "cc"}
	assert.Equal(t, []string{"a", "b", "c", "cc"}, Unique(ss))
}

func TestAny(t *testing.T) {
	ss := []string{"a", "b", "c", "bc"}
	hasTwoLetters := func(s string) bool {
		return len(s) == 2
	}
	assert.True(t, Any(ss, hasTwoLetters))
}

func TestAll(t *testing.T) {
	ss := []string{"a", "b", "c", "bc"}
	isNotEmpty := func(s string) bool {
		return s != ""
	}
	assert.True(t, All(ss, isNotEmpty))
}

type request struct {
	Id      int
	Latency time.Duration
	Done    chan struct{}
}

type response struct {
	Id int
}

var taskPool channel.Channel

func Service1(req *request) {
	taskPool.Input(req)
	return
}

func Service2(req *request) (*response, error) {
	if req.Latency > 0 {
		time.Sleep(req.Latency)
	}
	return &response{Id: req.Id}, nil
}

func TestNetworkIsolationOrDownstreamBlock(t *testing.T) {
	taskPool = channel.New(
		channel.WithNonBlock(),
		channel.WithTimeout(time.Millisecond*10),
	)
	defer taskPool.Close()
	var responded int32
	go func() {
		// task worker
		for task := range taskPool.Output() {
			req := task.(*request)
			done := make(chan struct{})
			go func() {
				_, _ = Service2(req)
				close(done)
			}()
			select {
			case <-time.After(time.Millisecond * 100):
			case <-done:
				atomic.AddInt32(&responded, 1)
			}
		}
	}()

	start := time.Now()
	for i := 1; i <= 100; i++ {
		req := &request{Id: i}
		if i > 50 && i <= 60 { // suddenly have network issue for 10 requests
			req.Latency = time.Hour
		}
		Service1(req)
	}
	cost := time.Now().Sub(start)
	assert.True(t, cost < time.Millisecond*10) // Service1 should not block
	time.Sleep(time.Millisecond * 1500)        // wait all tasks finished
	assert.Equal(t, int32(50), responded)      // 50 success and 10 timeout and 40 discard
}

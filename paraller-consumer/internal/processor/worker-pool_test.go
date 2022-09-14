package processor

import (
	"fmt"
	"testing"
	"time"
)

type PoolWorker1 struct {
	Name string
}

func (pW *PoolWorker1) DoWork(wR int) {
	time.Sleep(time.Second)
	fmt.Printf("Data: %s", pW.Name)
}

func TestBase(t *testing.T) {
	workerPool := New(100, 10000)

	for i := 0; i < 10000; i++ {
		workerPool.PostWork("routing", &PoolWorker1{Name: fmt.Sprintf("my namne %d", i)})
	}

}

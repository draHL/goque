package goque_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/draHL/goque"
)

func funcsleep(wait int) {
	fmt.Printf("running function %d \n", wait)
	time.Sleep(1 * time.Second)
}

func TestQue(t *testing.T) {
	maxQue := 10
	maxGoroutine := 0

	que, err := goque.NewJobQueue(maxQue, maxGoroutine)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	que.Start()
	defer que.Stop()

	for i := 0; i < 5; i++ {
		t.Logf("add job %d \n", i)
		j := i
		que.Add(func() { funcsleep(j) })
	}

}

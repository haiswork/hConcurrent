package hconcurrent

import (
	"sync"
	"testing"
	"time"
)

func TestConcurrent(t *testing.T) {
	inputChanSize := 4
	doFuncCount := 4
	listCount := 3
	inputs := []int{1, 2, 3, 4, 5}
	m := map[int]int{}
	l := new(sync.Mutex)

	c := NewConcurrent2(
		inputChanSize, doFuncCount, testDo,
		inputChanSize, doFuncCount, testDo,
		inputChanSize, doFuncCount, testDo,
		inputChanSize, doFuncCount, func(i interface{}) interface{} {
			l.Lock()
			defer l.Unlock()
			n := i.(int)
			m[n] = n
			return n
		},
	)

	c.Run()
	for i := 0; i < len(inputs); i++ {
		c.Input(inputs[i])
	}
	c.Stop()
	c.Destroy()

	for i := 0; i < len(inputs); i++ {
		n := inputs[i] + listCount
		if m[n] != n {
			t.Errorf("concurrent do error")
		}
	}
}

func TestConcurrentTimeout(t *testing.T) {
	c := NewConcurrent(1, 1, func(i interface{}) interface{} {
		time.Sleep(time.Second)
		return nil
	})
	c.Input(1)
	inputSuccess := c.InputWithTimeout(1, 100*time.Millisecond)
	if inputSuccess == true {
		t.Error("concurrent input with timeout error")
	}
}

func testDo(i interface{}) interface{} {
	n := i.(int)
	return n + 1
}

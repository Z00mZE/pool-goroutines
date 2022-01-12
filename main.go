package main

import (
	"fmt"
	"time"
)

func main() {
	pool := NewPool(4)
	for i := 0; i <= 16; i++ {
		pool.Exec(NewSimpleTask(fmt.Sprintf("%d", i)))
	}
	pool.Wait()
	pool.Close()

}

type DummyTask struct {
	message string
}

func NewSimpleTask(message string) *DummyTask {
	return &DummyTask{
		message: message,
	}
}
func (t *DummyTask) Execute() {
	time.Sleep(time.Millisecond * 300)
	fmt.Println(t.message)
}

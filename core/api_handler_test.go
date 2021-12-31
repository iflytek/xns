package core

import (
	"container/heap"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestFFF(t *testing.T) {

	tm := newTimer()
	tm.AddFunc(1*time.Second, func() {
		fmt.Println("done")
	})

	time.Sleep(2 *time.Second)
}

type timefunc struct {
	fun func()
	deadline time.Time
}

func newTimer()*timer{
	t := &timer{
		tasks: nil,
		mu:    sync.NewCond(&sync.Mutex{}),
	}
	t.runTimer()
	return t
}

type timer struct {
	tasks []*timefunc
	mu *sync.Cond
}

func (t *timer) Len() int {
	l:= len(t.tasks)
	return l
}

func (t *timer) Less(i, j int) bool {
	return t.tasks[i].deadline.Sub(t.tasks[j].deadline) < 0
}

func (t *timer) Swap(i, j int) {
	t.tasks[i],t.tasks[j] = t.tasks[j],t.tasks[i]
}

func (t *timer) Push(x interface{}) {
	t.tasks = append(t.tasks,x.(*timefunc))
}

func (t *timer) Pop() interface{} {
	v := t.tasks[len(t.tasks)-1]
	t.tasks = t.tasks[:len(t.tasks)-1]
	return v
}

func (t *timer)AddFunc(duration time.Duration,f func()){
	t.mu.L.Lock()
	heap.Push(t,&timefunc{
		fun:      f,
		deadline: time.Now().Add(duration),
	})
	t.mu.Signal()
	t.mu.L.Unlock()
}

func (t *timer)runTimer(){
	go func() {
		for{
			t.mu.L.Lock()
			if t.Len() == 0{
				t.mu.L.Unlock()
				time.Sleep(1*time.Second)
				continue
			}

		}
	}()
}

package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded     = errors.New("errors limit exceeded")
	ErrWrongNumberOfGoroutines = errors.New("wrong number of goroutines")
	ErrEmptyTaskList           = errors.New("empty task list")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return ErrEmptyTaskList
	}
	if n <= 0 {
		return ErrWrongNumberOfGoroutines
	}
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	wg := &sync.WaitGroup{}
	wg.Add(n)

	var errCount int32
	tasksCh := make(chan Task)

	for i := 0; i < n; i++ {
		go worker(wg, tasksCh, &errCount)
	}

	for _, t := range tasks {
		if atomic.LoadInt32(&errCount) >= int32(m) {
			break
		}
		tasksCh <- t
	}
	close(tasksCh)

	wg.Wait()
	if atomic.LoadInt32(&errCount) >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func worker(group *sync.WaitGroup, tasksCh chan Task, errCount *int32) {
	defer group.Done()
	for task := range tasksCh {
		if task() != nil {
			atomic.AddInt32(errCount, 1)
		}
	}
}

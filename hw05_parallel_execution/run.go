package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	ch := make(chan Task)
	wg := &sync.WaitGroup{}
	var errCount int32

	go func() {
		defer close(ch)
		for _, t := range tasks {
			ch <- t
		}
	}()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go doWork(ch, wg, m, &errCount)
	}
	wg.Wait()

	if int(atomic.LoadInt32(&errCount)) >= m && m > 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func doWork(ch chan Task, wg *sync.WaitGroup, maxErr int, errCount *int32) {
	for {
		isErrorLimitExceeded := int(atomic.LoadInt32(errCount)) >= maxErr

		// Desable error limit validation for m <= 0
		if maxErr <= 0 {
			isErrorLimitExceeded = false
		}

		chLen := len(ch)
		task, isOpen := <-ch

		// Execute task
		if !isErrorLimitExceeded && task != nil {
			err := task()
			if err != nil {
				atomic.AddInt32(errCount, 1)
			}
		}
		if chLen == 0 && !isOpen && task == nil {
			break
		}
	}
	wg.Done()
}

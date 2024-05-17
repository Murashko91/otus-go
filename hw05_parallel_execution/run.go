package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type workContext struct {
	mu       *sync.Mutex
	errCount int32
	wg       *sync.WaitGroup
	ch       chan Task
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	ch := make(chan Task)
	wg := &sync.WaitGroup{}
	wc := &workContext{mu: &sync.Mutex{}, wg: wg, ch: ch}
	errorLimitEnabled := m > 0

	go func() {
		defer close(ch)
		for _, t := range tasks {
			ch <- t
		}
	}()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go doWork(wc, m, errorLimitEnabled)
	}
	wg.Wait()

	if int(wc.errCount) >= m && errorLimitEnabled {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func doWork(wc *workContext, maxErr int, errorLimitEnabled bool) {
	for {
		wc.mu.Lock()
		isErrorLimitExceeded := int(wc.errCount) >= maxErr
		wc.mu.Unlock()

		// Desable error limit validation for m <= 0
		if !errorLimitEnabled {
			isErrorLimitExceeded = false
		}

		chLen := len(wc.ch)
		task, isOpen := <-wc.ch

		// Execute task
		if !isErrorLimitExceeded && task != nil {
			err := task()
			if err != nil {
				atomic.AddInt32(&wc.errCount, 1)
			}
		}

		if chLen == 0 && !isOpen && task == nil {
			break
		}
	}
	wc.wg.Done()
}

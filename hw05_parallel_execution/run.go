package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type workContext struct {
	tasks     []Task
	mu        *sync.Mutex
	errorsNum int
	tNumber   int
	wg        *sync.WaitGroup
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := &sync.WaitGroup{}
	wc := &workContext{tasks: tasks, mu: &sync.Mutex{}, wg: wg}
	errorLimitEnabled := m > 0

	for i := 0; i < n; i++ {
		wg.Add(1)
		go doWork(wc, m, errorLimitEnabled)
	}
	wg.Wait()

	if wc.errorsNum >= m && errorLimitEnabled {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func doWork(wc *workContext, maxErr int, errorLimitEnabled bool) {
	for {
		wc.mu.Lock()

		errLimit := wc.errorsNum >= maxErr

		// Desable error limit validation for m <= 0
		if !errorLimitEnabled {
			errLimit = false
		}

		taskIndex := wc.tNumber
		tasksCompleated := taskIndex > len(wc.tasks)-1

		wc.tNumber++

		if errLimit || tasksCompleated {
			wc.mu.Unlock()
			break
		}
		wc.mu.Unlock()

		// Execute task
		err := wc.tasks[taskIndex]()

		if err != nil && errorLimitEnabled {
			wc.mu.Lock()
			wc.errorsNum++
			wc.mu.Unlock()
		}
	}

	wc.wg.Done()
}

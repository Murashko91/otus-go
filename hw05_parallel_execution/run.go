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
	for i := 0; i < n; i++ {

		wg.Add(1)
		go doWork(wc, m)
	}
	wg.Wait()

	if wc.errorsNum >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func doWork(wc *workContext, m int) {

	for {

		wc.mu.Lock()

		errLimit := wc.errorsNum >= m
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

		if err != nil {
			wc.mu.Lock()
			wc.errorsNum++
			wc.mu.Unlock()
		}

	}

	wc.wg.Done()

}

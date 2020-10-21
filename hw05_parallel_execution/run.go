package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	if n <= 0 {
		n = 1
	}
	if m <= 0 {
		m = 0
	}

	mCasted := int32(m)
	taskCh := make(chan Task)
	quitCh := make(chan struct{})
	var errorsCounter int32 = 0

	wg := sync.WaitGroup{}
	worker := func(taskCh <-chan Task, quitCh <-chan struct{}) {
		defer wg.Done()
		for {
			select {
			case <-quitCh:
				return
			case task := <-taskCh:
				err := task()
				if err != nil {
					atomic.AddInt32(&errorsCounter, 1)
				}
			}
		}
	}

	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker(taskCh, quitCh)
	}

	for _, task := range tasks {
		if mCasted != 0 && atomic.LoadInt32(&errorsCounter) >= mCasted {
			break
		}
		taskCh <- task
	}
	close(quitCh)

	wg.Wait()
	close(taskCh)

	if mCasted != 0 && errorsCounter >= mCasted {
		return ErrErrorsLimitExceeded
	}

	return nil
}

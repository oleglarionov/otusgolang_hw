package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	if n <= 0 {
		n = 1
	}

	errorsChBufferSize := m
	if errorsChBufferSize <= 0 {
		errorsChBufferSize = len(tasks)
	}
	errorsCh := make(chan error, errorsChBufferSize)
	errorsCounter := 0

	successCh := make(chan struct{}, n)

	goroutinesCounter := 0
	wg := sync.WaitGroup{}
	i := 0
	tasksLen := len(tasks)
	completedTasks := 0
	var resultErr error
	for {
		if goroutinesCounter == n || i == tasksLen {
			select {
			case <-errorsCh:
				errorsCounter++
				completedTasks++
				if errorsCounter == m {
					resultErr = ErrErrorsLimitExceeded
				}
				goroutinesCounter--
			case <-successCh:
				completedTasks++
				goroutinesCounter--
			}

			if resultErr != nil || completedTasks == tasksLen {
				break
			}
		}

		if i < tasksLen {
			goroutinesCounter++

			wg.Add(1)
			go func(task Task, errorsCh chan<- error, successCh chan<- struct{}) {
				defer wg.Done()

				err := task()
				if err != nil {
					errorsCh <- err
				} else {
					successCh <- struct{}{}
				}
			}(tasks[i], errorsCh, successCh)

			i++
		}
	}
	wg.Wait()

	return resultErr
}

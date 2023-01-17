package hw05parallelexecution

import (
	"errors"
	"sync"
)

type Task func() error

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

var ErrNoThreadsForExecute = errors.New("errors no threads for execute")

func Run(tasks []Task, n, m int) (err error) {
	if n <= 0 {
		err = ErrNoThreadsForExecute
		return
	}
	wg := &sync.WaitGroup{}
	if m < 0 {
		m = len(tasks) + 1
	}
	chTask := make(chan Task)
	chErr := make(chan error, n+1)
	defer func() {
		close(chTask)
		wg.Wait()
		close(chErr)
		for range chErr {
			if m--; m == 0 {
				err = ErrErrorsLimitExceeded
			}
		}
	}()
	if n > len(tasks) {
		n = len(tasks)
	}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range chTask {
				if err := t(); err != nil {
					chErr <- err
				}
			}
		}()
	}
	for _, task := range tasks {
		for {
			exit := false
			select {
			case <-chErr:
				if m--; m == 0 {
					return
				}
			default:
				exit = true
			}
			if exit {
				break
			}
		}
		chTask <- task
	}
	return
}

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
	if m < 0 {
		m = len(tasks) + 1
	}
	chTask := make(chan Task)
	chErr := make(chan error, n+1)
	if n > len(tasks) {
		n = len(tasks)
	}
	wg := &sync.WaitGroup{}
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
				m--
			default:
				exit = true
			}
			if exit {
				break
			}
		}
		if m <= 0 {
			break
		}
		chTask <- task
	}
	close(chTask)
	wg.Wait()
	close(chErr)
	for range chErr {
		m--
	}
	if m <= 0 {
		err = ErrErrorsLimitExceeded
	}
	return
}

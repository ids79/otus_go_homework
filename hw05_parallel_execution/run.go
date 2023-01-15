package hw05parallelexecution

import (
	"errors"
	"sync"
)

type Task func() error

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

var ErrNoThreadsForExecute = errors.New("errors no threads for execute")

func Run(tasks []Task, n, m int) (err error) {
	wg := &sync.WaitGroup{}
	if m < 0 {
		m = len(tasks) + 1
	}
	chTask := make(chan Task, 1)
	chErr := make(chan error, m)
	defer func() {
		wg.Wait()
		close(chErr)
		for {
			if _, ok := <-chErr; ok {
				m--
			} else {
				break
			}
		}
		if m <= 0 {
			err = ErrErrorsLimitExceeded
		}
	}()
	if n <= 0 {
		err = ErrNoThreadsForExecute
		return
	}
	if n > len(tasks) {
		n = len(tasks)
	}
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for t := range chTask {
				err := t()
				if err != nil {
					chErr <- err
				}
			}
		}()
		wg.Add(1)
	}
	for i := 0; i < len(tasks); {
		select {
		case <-chErr:
			m--
		default:
			if m == 0 {
				close(chTask)
				return
			}
			select {
			case chTask <- tasks[i]:
				i++
			default:
			}
		}
	}
	close(chTask)
	return
}

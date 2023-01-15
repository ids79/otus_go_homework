package hw05parallelexecution

import (
	"errors"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

var chTask chan Task

var chErr chan error

var working int32

func worcker(t Task) {
	defer atomic.AddInt32(&working, -1)
	err := t()
	if err != nil {
		chErr <- err
		return
	}
	for {
		if t, ok := <-chTask; ok {
			err = t()
			if err != nil {
				chErr <- err
				return
			}
		} else {
			return
		}
	}
}

func Run(tasks []Task, n, m int) (err error) {
	defer func() {
		for atomic.LoadInt32(&working) != 0 {
			select {
			case <-chErr:
				m--
			default:
			}
		}
		if m <= 0 {
			err = ErrErrorsLimitExceeded
		}
	}()
	chTask = make(chan Task, 1)
	chErr = make(chan error)
	working = 0
	if m < 0 {
		m = len(tasks) + 1
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
			if int(atomic.LoadInt32(&working)) < n {
				go worcker(tasks[i])
				atomic.AddInt32(&working, 1)
				i++
			} else {
				select {
				case chTask <- tasks[i]:
					i++
				default:
				}
			}
		}
	}
	if len(chTask) == 1 {
		go worcker(<-chTask)
		atomic.AddInt32(&working, 1)
	}
	close(chTask)
	return
}

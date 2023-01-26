package hw06pipelineexecution

import (
	"sync"
	"sync/atomic"
	"time"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

var Done int32

func ExecutePipeline(in In, done In, wg *sync.WaitGroup, stages ...Stage) Out {
	Done = 0
	for _, stage := range stages {
		in = stage(in)
	}
	if done != nil {
		go func() {
			tick := time.NewTicker(time.Millisecond * 10)
			for {
				<-tick.C
				if isClose(done) {
					atomic.AddInt32(&Done, 1)
					break
				}
			}
			tick.Stop()
		}()
	}
	return in
}

func isClose(done In) bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

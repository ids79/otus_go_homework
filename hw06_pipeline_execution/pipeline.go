package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
	Val struct {
		value interface{}
		done  In
	}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		wr := make(chan interface{})
		go func(in In) {
			defer close(wr)
			for v := range in {
				d := Val{
					value: v,
					done:  done,
				}
				select {
				case <-done:
					return
				case wr <- d:
				}
			}
		}(in)
		in = stage(wr)
	}
	return in
}

package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	stage := stages[0]
	result := make(Bi)

	if len(stages) > 1 {
		exResult := ExecutePipeline(stage(in), done, stages[1:]...)

		go func() {
			for {
				select {
				case <-done:
					close(result)
					return
				case i, isOpen := <-exResult:
					if !isOpen {
						close(result)
						return
					}
					result <- i
				}
			}
		}()
	} else {
		return stage(in)
	}
	return result
}

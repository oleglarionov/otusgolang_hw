package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	stageWrapper := func(in In, done In, stage Stage) Out {
		stageIn := make(Bi)
		out := stage(stageIn)

		go func() {
			for {
				select {
				case <-done:
					close(stageIn)
					return
				case v, ok := <-in:
					if !ok {
						close(stageIn)
						return
					}
					stageIn <- v
				}
			}
		}()

		return out
	}

	for _, stage := range stages {
		in = stageWrapper(in, done, stage)
	}

	return in
}

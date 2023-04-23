package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return nil
	}

	out := in
	for i := range stages {
		out = execStage(stages[i], out, done)
	}
	return out
}

func execStage(stage Stage, in In, done In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case data, ok := <-in:
				if ok {
					out <- data
				} else {
					return
				}
			}
		}
	}()
	return stage(out)
}

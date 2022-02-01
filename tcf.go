package noob

import (
	"runtime"
)

type Func struct {
	Try     func()
	Catch   func(e interface{}, frames *runtime.Frames)
	Finally func()
}

func TCFunc(run Func) {
	if run.Finally != nil {
		defer run.Finally()
	}

	defer func() {
		if r := recover(); r != nil {
			if run.Catch != nil {
				run.Catch(r, GetRuntimeFrames(4))
				return
			}

			panic(r)
		}
	}()

	run.Try()
}

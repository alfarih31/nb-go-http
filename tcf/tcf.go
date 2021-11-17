package tcf

import (
	apperr "github.com/alfarih31/nb-go-http/app_err"
	"runtime"
)

type Func struct {
	Try     func()
	Catch   func(e interface{}, frames *runtime.Frames)
	Finally func()
}

func TCFunc(ru Func) {
	if ru.Finally != nil {
		defer ru.Finally()
	}

	defer func() {
		if r := recover(); r != nil {
			if ru.Catch != nil {
				ru.Catch(r, apperr.GetRuntimeFrames(4))
				return
			}

			panic(r)
		}
	}()

	ru.Try()
}

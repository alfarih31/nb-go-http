package nbgohttp

import "runtime"

type FuncCtx struct {
	Try     func()
	Catch   func(e interface{})
	Finally func()
}

func Func(t FuncCtx) {
	if t.Finally != nil {
		defer t.Finally()
	}

	defer func() {
		if r := recover(); r != nil {
			if t.Catch != nil {
				t.Catch(r)
				return
			}

			panic(r)
		}
	}()

	runtime.StartTrace()
	t.Try()
}

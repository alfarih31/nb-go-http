package utils

type Run struct {
	Try     func()
	Catch   func(e interface{})
	Finally func()
}

func Func(ru Run) {
	if ru.Finally != nil {
		defer ru.Finally()
	}

	defer func() {
		if r := recover(); r != nil {
			if ru.Catch != nil {
				ru.Catch(r)
				return
			}

			panic(r)
		}
	}()

	ru.Try()
}

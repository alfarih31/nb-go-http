package tcf

type Func struct {
	Try     func()
	Catch   func(e interface{})
	Finally func()
}

func TCFunc(ru Func) {
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

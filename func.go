package nbgohttp

type Func struct {
	Try     func()
	Catch   func(e interface{})
	Finally func()
}

func FlowFunc(ru Func) {
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

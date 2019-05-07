package runtime

// ReallyCrash 程序真正挂掉
const ReallyCrash = true

// TODO 增加默认处理函数
var PanicHandlers = []func(interface{}){}

// HandleCrash 简单的捕获异常，注意通过defer调用.
func HandleCrash(additionalHandlers ...func(interface{})) {
	if r := recover(); r != nil {
		for _, fn := range PanicHandlers {
			fn(r)
		}
		for _, fn := range additionalHandlers {
			fn(r)
		}
		if ReallyCrash {
			// Actually proceed to panic.
			panic(r)
		}
	}
}

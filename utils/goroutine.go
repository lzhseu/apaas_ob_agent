package utils

import (
	"fmt"
	"runtime"
)

func Go(f func()) {
	go func() {
		defer PanicGuard()
		f()
	}()
}

func PanicGuard() {
	if r := recover(); r != nil {
		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]
		//logs.CtxError(ctx, "🍀env: %v, panic recover: %s\n%s", env.Env(), format.Any(r), buf)
		fmt.Printf("🍀env: %v, panic recover: %+v\n%s\n", "test", r, buf)
	}
}

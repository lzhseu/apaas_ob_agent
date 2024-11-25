package recovery

import (
	"fmt"
	"runtime"

	"github.com/lzhseu/apaas_ob_agent/inner/logs"
)

type OptionFunc func(g *Option)

type Option struct {
	recoverHandler func(r any)
}

func WithRecoverHandler(recoverHandler func(r any)) OptionFunc {
	return func(g *Option) {
		g.recoverHandler = recoverHandler
	}
}

func Go(f func(), opts ...OptionFunc) {
	go func() {
		defer PanicGuard(opts...)
		f()
	}()
}

func PanicGuard(opts ...OptionFunc) {
	opt := &Option{}
	for _, v := range opts {
		v(opt)
	}

	if r := recover(); r != nil {
		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]
		logs.Error(fmt.Sprintf("🍀panic recover: %+v\n%s\n", r, buf))

		if opt.recoverHandler != nil {
			opt.recoverHandler(r)
		}
	}
}

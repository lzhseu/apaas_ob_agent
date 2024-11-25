package inner

import (
	"github.com/lzhseu/apaas_ob_agent/inner/logs"
	"github.com/lzhseu/apaas_ob_agent/inner/metrics"
)

func MustInit() {
	metrics.MustInit()
	logs.MustInit()
}

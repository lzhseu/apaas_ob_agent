package service

import "github.com/lzhseu/apaas_ob_agent/service/metrics"

func MustInit() {
	metrics.MustInit()
}

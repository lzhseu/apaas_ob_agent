package service

import (
	"github.com/lzhseu/apaas_ob_agent/service/prometheus"
)

func MustInit() {
	prometheus.MustInit()
}

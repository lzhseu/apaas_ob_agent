package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/lzhseu/apaas_ob_agent/conf"
	feishu "github.com/lzhseu/apaas_ob_agent/feishu_event"
	"github.com/lzhseu/apaas_ob_agent/inner"
	"github.com/lzhseu/apaas_ob_agent/service"
)

func main() {
	conf.MustInit()
	inner.MustInit()
	feishu.MustInit()
	service.MustInit()

	http.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{Registry: prometheus.DefaultRegisterer}))

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", conf.GetConfig().HttpServerCfg.Host, conf.GetConfig().HttpServerCfg.Port), nil))
}

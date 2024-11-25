package main

import (
	"fmt"
	"log"
	"net/http"

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

	http.HandleFunc("/ping", Ping)
	http.HandleFunc("/metrics", PrometheusExporter)
	http.HandleFunc("/alert", AlertWebhook)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", conf.GetConfig().HttpServerCfg.Host, conf.GetConfig().HttpServerCfg.Port), nil))
}

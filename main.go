package main

import (
	"fmt"
	"net/http"

	"github.com/lzhseu/apaas_ob_agent/config"
	feishu "github.com/lzhseu/apaas_ob_agent/feishu_event"
	"github.com/lzhseu/apaas_ob_agent/inner"
	"github.com/lzhseu/apaas_ob_agent/inner/logs"
	"github.com/lzhseu/apaas_ob_agent/service"
)

func main() {
	config.MustInit()
	inner.MustInit()
	feishu.MustInit()
	service.MustInit()

	http.HandleFunc("/ping", Ping)
	http.HandleFunc("/metrics", PrometheusExporter)

	startHttpServer()
}

func startHttpServer() {
	logs.Info("http server start...")

	rootURL := config.GetConfig().HttpServerCfg.Host
	if port := config.GetConfig().HttpServerCfg.Port; port != "" {
		rootURL = fmt.Sprintf("%v:%v", rootURL, port)
	}

	err := http.ListenAndServe(rootURL, nil)
	if err != nil {
		logs.Error(fmt.Sprintf("http server exit unexpectedly. err = %v", err))
	}
	logs.Info("http server exit")
}

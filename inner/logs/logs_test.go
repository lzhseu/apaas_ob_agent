package logs

import (
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/lzhseu/apaas_ob_agent/conf"
)

func TestLog(t *testing.T) {
	gomonkey.ApplyFuncReturn(conf.GetConfig, &conf.Config{
		InnerLogsCfg: &conf.LogsCfg{
			Console: &conf.LogConsoleCfg{Enable: true},
			File:    &conf.LogFileCfg{Enable: true, FileName: "test.log"},
			Loki:    &conf.LogLokiCfg{Enable: true, Schema: "http", Host: "10.37.107.200", Port: "3100", Labels: map[string]string{"service_name": "apaas_agent"}}},
	})
	MustInit()
	Info("this is a info log", "foo", "bar")
	Warn("this is a warn log", "foo", "bar")
	Error("this is a error log", "foo", "bar")
	time.Sleep(time.Second * 3)
}

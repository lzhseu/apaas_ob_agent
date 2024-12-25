package feishu_event

import (
	"context"
	"fmt"
	"time"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"

	"github.com/lzhseu/apaas_ob_agent/config"
	"github.com/lzhseu/apaas_ob_agent/inner/logs"
)

func MustInit() {
	clientID := config.GetConfig().FeishuAppCfg.AppID
	clientSecret := config.GetConfig().FeishuAppCfg.AppSecret
	logLevel := larkcore.LogLevelInfo
	if config.GetConfig().InnerLogsCfg.LogLevel != nil {
		logLevel = logLevelFromStr[*(config.GetConfig().InnerLogsCfg.LogLevel)]
	}

	// 注册「事件-事件处理器」
	eventHandler := dispatcher.NewEventDispatcher("", "").
		OnCustomizedEvent(EventTypeMetricReported, NewFeishuEventHandler(EventTypeMetricReported, NewMetricsBizHandler))

	cli := larkws.NewClient(clientID, clientSecret,
		larkws.WithEventHandler(eventHandler),
		larkws.WithLogLevel(logLevel),
	)

	errChan := make(chan error)
	timeout := time.After(time.Second * 3)
	go func() {
		if err := cli.Start(context.Background()); err != nil {
			errChan <- err
			return
		}
	}()

	select {
	case <-timeout:
		logs.Info("feishu event listener started")
	case err := <-errChan:
		logs.Error(fmt.Sprintf("feishu event listener start err: %v", err))
		panic(err)
	}
}

var (
	logLevelFromStr = map[string]larkcore.LogLevel{
		"debug": larkcore.LogLevelDebug,
		"info":  larkcore.LogLevelInfo,
		"warn":  larkcore.LogLevelWarn,
		"error": larkcore.LogLevelError,
	}
)

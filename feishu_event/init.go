package feishu_event

import (
	"context"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"

	"github.com/lzhseu/apaas_ob_agent/conf"
)

func MustInit() {
	clientID := conf.GetConfig().FeishuAppCfg.AppID
	clientSecret := conf.GetConfig().FeishuAppCfg.AppSecret
	logLevel := larkcore.LogLevelInfo
	if conf.GetConfig().InnerLogsCfg.LogLevel != nil {
		logLevel = logLevelFromStr[*(conf.GetConfig().InnerLogsCfg.LogLevel)]
	}

	// 注册「事件-事件处理器」
	eventHandler := dispatcher.NewEventDispatcher("", "").
		OnCustomizedEvent(EventTypeMetricReported, NewFeishuEventHandler(EventTypeMetricReported, NewMetricsBizHandler))

	cli := larkws.NewClient(clientID, clientSecret,
		larkws.WithEventHandler(eventHandler),
		larkws.WithLogLevel(logLevel),
	)

	go func() {
		if err := cli.Start(context.Background()); err != nil {
			panic(err)
		}
	}()
}

var (
	logLevelFromStr = map[string]larkcore.LogLevel{
		"debug": larkcore.LogLevelDebug,
		"info":  larkcore.LogLevelInfo,
		"warn":  larkcore.LogLevelWarn,
		"error": larkcore.LogLevelError,
	}
)

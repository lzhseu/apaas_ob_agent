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

	// 注册「事件-事件处理器」
	eventHandler := dispatcher.NewEventDispatcher("", "").
		OnCustomizedEvent(EventTypeMetricReported, NewFeishuEventHandler(EventTypeMetricReported, NewMetricsBizHandler))

	cli := larkws.NewClient(clientID, clientSecret,
		larkws.WithEventHandler(eventHandler),
		larkws.WithLogLevel(larkcore.LogLevelDebug),
		larkws.WithDomain("https://open.feishu-boe.cn")) // todo：boe调试用，正式上线删除

	go func() {
		if err := cli.Start(context.Background()); err != nil {
			panic(err)
		}
	}()
}

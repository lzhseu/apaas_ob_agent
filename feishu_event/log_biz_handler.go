package feishu_event

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/lzhseu/apaas_ob_agent/config"
	"github.com/lzhseu/apaas_ob_agent/inner/logs"
	"github.com/lzhseu/apaas_ob_agent/pkg/logger"
)

var (
	logLoki *logger.Logger
)

func mustInitLogLogger() {
	rootURL := config.GetConfig().InnerLogsCfg.Loki.RootURL
	labels := config.GetConfig().InnerLogsCfg.Loki.Labels
	labels["data_type"] = "log"
	logLoki = logger.NewLokiLogger(rootURL, labels)
}

type LogBizHandler struct {
	Data *LogEventData
}

func (e *LogBizHandler) Unmarshal(ctx context.Context, packet *FeishuEventPacket) error {
	e.Data = &LogEventData{}
	if err := sonic.Unmarshal(packet.Event, e.Data); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (e *LogBizHandler) Validate(ctx context.Context, packet *FeishuEventPacket) error {
	validate := validator.New()
	err := validate.Struct(e.Data)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (e *LogBizHandler) Handle(ctx context.Context, packet *FeishuEventPacket) error {
	for _, v := range e.Data.Logs {
		if err := e.handleLog(v); err != nil {
			logs.Error(fmt.Sprintf("[LogBizHandler] handle log error: %v", err), v.genInnerLogTags()...)
			continue
		}
	}
	return nil
}

func (e *LogBizHandler) handleLog(v *Log) (err error) {
	str, err := sonic.MarshalString(v)
	if err != nil {
		logLoki.Error(fmt.Sprintf("marshal log error: %v", err))
		return nil
	}
	logLoki.Info(str)
	return nil
}

// todo
func (m *Log) genInnerLogTags() []any {
	return []any{
		"event_type", EventTypeLogReported,
	}
}

func NewLogBizHandler() BizHandlerIface {
	return &LogBizHandler{}
}

type LogEventData struct {
	Logs []*Log `json:"logs"`
}

type Log struct {
	Content    string            `json:"content,omitempty"`
	Level      string            `json:"level,omitempty"`
	Timestamp  int64             `json:"timestamp,omitempty"` // 单位：毫秒
	BizEvent   *Event            `json:"event,omitempty"`
	TraceID    string            `json:"trace_id,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"` // 附件属性，包括应用相关属性（如：tenantID，namespace 等）
}

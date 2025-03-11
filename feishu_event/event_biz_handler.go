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
	eventLoki *logger.Logger
)

func mustInitEventLogger() {
	rootURL := config.GetConfig().InnerLogsCfg.Loki.RootURL
	labels := config.GetConfig().InnerLogsCfg.Loki.Labels
	labels["data_type"] = "event"
	eventLoki = logger.NewLokiLogger(rootURL, labels)
}

type EventBizHandler struct {
	Data *BizEventData
}

func (e *EventBizHandler) Unmarshal(ctx context.Context, packet *FeishuEventPacket) error {
	e.Data = &BizEventData{}
	if err := sonic.Unmarshal(packet.Event, e.Data); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (e *EventBizHandler) Validate(ctx context.Context, packet *FeishuEventPacket) error {
	validate := validator.New()
	err := validate.Struct(e.Data)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (e *EventBizHandler) Handle(ctx context.Context, packet *FeishuEventPacket) error {
	for _, event := range e.Data.Events {
		if err := e.handleEvent(event); err != nil {
			logs.Error(fmt.Sprintf("[EventBizHandler] handle event error: %v", err), event.genInnerLogTags()...)
			continue
		}
	}
	return nil
}

func (e *EventBizHandler) handleEvent(event *Event) (err error) {
	eventLoki.Info(fmt.Sprintf("%#v", event))
	return nil
}

// todo
func (m *Event) genInnerLogTags() []any {
	return []any{
		"event_type", EventTypeEventReported,
		"id", m.ID,
		"type", m.Type,
	}
}

func NewEventBizHandler() BizHandlerIface {
	return &EventBizHandler{}
}

type BizEventData struct {
	Events []*Event `json:"events"`
}

type Event struct {
	ID             string            `json:"id,omitempty"`   // 观测事件ID
	Type           string            `json:"type,omitempty"` // 观测事件类型
	TraceID        string            `json:"trace_id,omitempty"`
	StartTimestamp int64             `json:"start_timestamp,omitempty"`
	EndTimestamp   int64             `json:"end_timestamp,omitempty"`
	IsFinished     bool              `json:"is_finished,omitempty"`
	Detail         string            `json:"detail,omitempty"`     // 观测事件详情
	Attributes     map[string]string `json:"attributes,omitempty"` // 附加属性，包括应用相关属性（如：tenantID，namespace 等）
}

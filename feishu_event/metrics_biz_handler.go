package feishu_event

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/lzhseu/apaas_ob_agent/service/metrics"
)

type MetricsBizHandler struct {
	Data *MetricEventData
}

func NewMetricsBizHandler() BizHandlerIface {
	return &MetricsBizHandler{}
}

func (m *MetricsBizHandler) Unmarshal(ctx context.Context, packet *FeishuEventPacket) error {
	m.Data = &MetricEventData{}
	if err := sonic.Unmarshal(packet.Event, m.Data); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (m *MetricsBizHandler) Validate(ctx context.Context, packet *FeishuEventPacket) error {
	validate := validator.New()
	err := validate.Struct(m.Data)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (m *MetricsBizHandler) Handle(ctx context.Context, packet *FeishuEventPacket) error {
	for _, metric := range m.Data.Metrics {
		fmt.Printf("\n[lzh] metric=%#v\n", metric)
		labelNames := make([]string, 0, len(metric.Attributes))
		for key := range metric.Attributes {
			labelNames = append(labelNames, key)
		}

		collector, err := metrics.GetOrCreateCollector(metric.Name, metric.Type, labelNames)
		if err != nil {
			fmt.Printf("\n[lzh] GetOrCreateCollector err: %v\n", err)
			// todo: log
		}

		fmt.Printf("\n[lzh] GetOrCreateCollector: %#v\n", collector)

		collector.Collect(metric.Attributes, metric.Value)
	}

	return nil
}

type MetricEventData struct {
	Metrics []*Metric `json:"metrics"`
}

type Metric struct {
	Name       string `validate:"required"`
	Type       string `validate:"required"` // 指标类型
	Value      float64
	Attributes map[string]string
	Timestamp  int64
}

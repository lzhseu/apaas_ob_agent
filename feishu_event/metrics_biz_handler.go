package feishu_event

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/lzhseu/apaas_ob_agent/inner/logs"
	innermetrics "github.com/lzhseu/apaas_ob_agent/inner/metrics"
	"github.com/lzhseu/apaas_ob_agent/service/prometheus"
)

const (
	metricExpiredInterval = 15 * 60 * 1000 // 15min
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
		if err := m.handlerMetric(metric); err != nil {
			logs.Error(fmt.Sprintf("[MetricsBizHandler] handle metric error: %v", err), metric.genInnerLogTags()...)
			continue
		}
	}
	return nil
}

func (m *MetricsBizHandler) handlerMetric(metric *Metric) (err error) {
	startAt := time.Now()
	var isDiscard bool

	defer func() {
		// agent 自身指标收集
		cost := time.Since(startAt).Milliseconds()
		attr := make(map[string]string)
		if metric.Attributes != nil {
			attr = metric.Attributes
		}

		labelVal := []string{
			attr["tenant_id"],
			attr["tenant_type"],
			attr["namespace"],
			metric.Name,
			metric.Type,
		}

		isErr := "false"
		errMsg := "-"
		if err != nil {
			isErr = "true"
			errMsg = errors.Cause(err).Error()
		} else {
			innermetrics.MetricEventHandleDurationMsHistogram.WithLabelValues(labelVal...).Observe(float64(cost))
			innermetrics.MetricEventHandleDurationMsSummary.WithLabelValues(labelVal...).Observe(float64(cost))
		}
		labelVal = append(labelVal, isErr, errMsg, fmt.Sprintf("%v", isDiscard))
		innermetrics.MetricEventHandleTotal.WithLabelValues(labelVal...).Inc()
	}()

	// 过期数据（15分钟前）不处理
	if metric.Timestamp > 0 && time.Now().UnixMilli()-metric.Timestamp > metricExpiredInterval {
		isDiscard = true
		logs.Warn(fmt.Sprintf("[handlerMetric] metric is discard as timestamp is invalid. metric timestamp: %v, now: %v", metric.Timestamp, time.Now().UnixMilli()), metric.genInnerLogTags()...)
		return nil
	}

	labelNames := make([]string, 0, len(metric.Attributes))
	for key := range metric.Attributes {
		labelNames = append(labelNames, key)
	}

	// 获取 collector
	collector, err := prometheus.GetOrCreateCollector(metric.Name, metric.Type, labelNames)
	if err != nil {
		return err
	}

	// 收集指标
	if err = collector.Collect(metric.Attributes, metric.Value); err != nil {
		return err
	}

	return nil
}

func (m *Metric) genInnerLogTags() []any {
	return []any{
		"event_type", EventTypeMetricReported,
		"metric_name", m.Name,
		"metric_type", m.Type,
	}
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

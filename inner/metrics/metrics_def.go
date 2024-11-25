package metrics

import "github.com/prometheus/client_golang/prometheus"

// agent 自身的监控指标
var (
	PanicTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "agent",
		Name:      "panic",
	}, []string{"scene"})

	FeishuEventReceiveTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "agent",
		Name:      "feishu_event_receive_total",
	}, []string{"feishu_event_name", "schema", "is_error", "error_msg"})

	FeishuEventReceiveDurationMsHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "agent",
		Name:      "feishu_event_receive_duration_ms_histogram",
	}, []string{"feishu_event_name", "schema"})

	FeishuEventReceiveDurationMsSummary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "agent",
		Name:      "feishu_event_receive_duration_ms_summary",
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001},
	}, []string{"feishu_event_name", "schema"})

	MetricEventHandleTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "agent",
		Name:      "metric_event_handle_total",
	}, []string{"tenant_id", "tenant_type", "namespace", "metric_name", "metric_type", "is_error", "error_msg", "is_discard"})

	MetricEventHandleDurationMsHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "agent",
		Name:      "metric_event_handle_duration_ms_histogram",
	}, []string{"tenant_id", "tenant_type", "namespace", "metric_name", "metric_type"})

	MetricEventHandleDurationMsSummary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "agent",
		Name:      "metric_event_handle_duration_ms_summary",
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001},
	}, []string{"tenant_id", "tenant_type", "namespace", "metric_name", "metric_type"})
)

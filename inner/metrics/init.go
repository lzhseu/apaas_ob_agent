package metrics

import "github.com/prometheus/client_golang/prometheus"

func MustInit() {
	prometheus.MustRegister(
		PanicTotal,
		FeishuEventReceiveTotal,
		FeishuEventReceiveDurationMsHistogram,
		FeishuEventReceiveDurationMsSummary,
		MetricEventHandleTotal,
		MetricEventHandleDurationMsHistogram,
		MetricEventHandleDurationMsSummary,
	)
}

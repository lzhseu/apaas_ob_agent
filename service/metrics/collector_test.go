package metrics

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/lzhseu/apaas_ob_agent/conf"
	"github.com/prometheus/client_golang/prometheus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCollector(t *testing.T) {
	Convey("", t, func() {

	})
}

func TestGenOpt(t *testing.T) {
	Convey("all", t, func() {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		patches.ApplyFuncReturn(conf.GetConfig,
			&conf.Config{
				PrometheusCfg: map[string]*conf.PrometheusCfg{
					"test_histogram_metric": {
						Name:    "test_histogram_metric",
						Type:    "histogram",
						Buckets: []float64{1, 2, 3},
					},
					"test_histogram_metric2": {
						Name: "test_histogram_metric2",
						Type: "histogram",
					},
					"test_summary_metric": {
						Name: "test_summary_metric",
						Type: "summary",
						Objectives: map[float64]float64{
							0.5:  0.05,
							0.9:  0.01,
							0.99: 0.001,
						},
					},
				},
			},
		)

		val := genOpts("no_exist_name", "counter")
		counterOpts, ok := val.(*prometheus.CounterOpts)
		So(ok, ShouldBeTrue)
		So(*counterOpts, ShouldEqual, prometheus.CounterOpts{Name: "no_exist_name"})

		val = genOpts("test_histogram_metric2", "histogram")
		hisOpts, ok := val.(*prometheus.HistogramOpts)
		So(ok, ShouldBeTrue)
		So(*hisOpts, ShouldEqual, prometheus.HistogramOpts{Name: "test_histogram_metric2"})

		val = genOpts("test_histogram_metric", "histogram")
		hisOpts, ok = val.(*prometheus.HistogramOpts)
		So(ok, ShouldBeTrue)
		So(hisOpts.Buckets, ShouldNotBeEmpty)
		t.Logf("")

		val = genOpts("test_summary_metric", "summary")
		summaryOpts, ok := val.(*prometheus.SummaryOpts)
		So(ok, ShouldBeTrue)
		So(summaryOpts.Objectives, ShouldNotBeEmpty)
	})
}

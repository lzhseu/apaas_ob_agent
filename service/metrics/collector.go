package metrics

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/lzhseu/apaas_ob_agent/utils"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/lzhseu/apaas_ob_agent/conf"
)

var (
	collectors map[string]*Collector
	mu         sync.RWMutex
)

type Collector struct {
	Name string // 指标名称
	Type string // 指标类型
	prom prometheus.Collector
}

func (c *Collector) Collect(tags map[string]string, value float64) {
	defer utils.PanicGuard()

	switch c.prom.(type) {
	case *prometheus.CounterVec:
		c.prom.(*prometheus.CounterVec).With(tags).Add(value)
	case *prometheus.GaugeVec:
		c.prom.(*prometheus.GaugeVec).With(tags).Set(value)
	case *prometheus.HistogramVec:
		c.prom.(*prometheus.HistogramVec).With(tags).Observe(value)
	case *prometheus.SummaryVec:
		c.prom.(*prometheus.SummaryVec).With(tags).Observe(value)
	default:
		// todo: log
	}
}

func MustInit() {
	collectors = make(map[string]*Collector)

	// 初始化配置文件中的指标
	// 配置文件中主要是有一些特殊配置，如果没有配置在文件中，则运行时也会动态创建 collector，但此时使用的是 Prometheus SDK 的默认配置
	for _, cfg := range conf.GetConfig().PrometheusCfg {
		collector, err := createCollector(cfg)
		if err != nil {
			panic(err)
		}
		collectors[collector.Name] = collector
	}
}

func GetOrCreateCollector(name, typ string, labelNames []string) (*Collector, error) {
	mu.RLock()
	if c, ok := collectors[name]; ok {
		mu.RUnlock()
		return c, nil
	}
	mu.RUnlock()

	// 创建新的 collector 前需要加互斥锁
	mu.Lock()
	defer mu.Unlock()

	// 再次检查是否已经创建了 collector
	if c, ok := collectors[name]; ok {
		return c, nil
	}

	// 创建新的 collector
	c, err := createCollector(&conf.PrometheusCfg{
		Name:       name,
		Type:       typ,
		LabelNames: labelNames,
	})
	if err != nil {
		return nil, err
	}

	collectors[c.Name] = c

	return c, nil
}

func createCollector(cfg *conf.PrometheusCfg) (*Collector, error) {
	var collector *Collector

	switch cfg.Type {
	case MTypeCounter:
		collector = createCounterCollector(cfg.Name, cfg.LabelNames)
	case MTypeGauge:
		collector = createGaugeCollector(cfg.Name, cfg.LabelNames)
	case MTypeHistogram:
		collector = createHistogramCollector(cfg.Name, cfg.LabelNames)
	case MTypeSummary:
		collector = createSummaryCollector(cfg.Name, cfg.LabelNames)
	default:
		return nil, errors.Errorf("[createCollector] invalid metric type: %s", cfg.Type)
	}

	fmt.Printf("\n[lzh] create collector=%#v\n", collector)
	prometheus.MustRegister(collector.prom)

	return collector, nil
}

func createCounterCollector(name string, labelNames []string) *Collector {
	opts := genOpts(name, MTypeCounter).(*prometheus.CounterOpts)
	return &Collector{
		Name: name,
		Type: MTypeCounter,
		prom: prometheus.NewCounterVec(*opts, labelNames),
	}
}

func createGaugeCollector(name string, labelNames []string) *Collector {
	opts := genOpts(name, MTypeGauge).(*prometheus.GaugeOpts)
	return &Collector{
		Name: name,
		Type: MTypeGauge,
		prom: prometheus.NewGaugeVec(*opts, labelNames),
	}
}

func createHistogramCollector(name string, labelNames []string) *Collector {
	opts := genOpts(name, MTypeHistogram).(*prometheus.HistogramOpts)
	return &Collector{
		Name: name,
		Type: MTypeHistogram,
		prom: prometheus.NewHistogramVec(*opts, labelNames),
	}
}

func createSummaryCollector(name string, labelNames []string) *Collector {
	opts := genOpts(name, MTypeSummary).(*prometheus.SummaryOpts)
	return &Collector{
		Name: name,
		Type: MTypeSummary,
		prom: prometheus.NewSummaryVec(*opts, labelNames),
	}
}

func genOpts(name, typ string) (opts any) {
	opts = defaultOpts(name, typ)

	// 如果配置文件中有，优先使用配置文件的参数，没有则使用默认参数
	cfg, ok := conf.GetConfig().PrometheusCfg[name]
	if !ok {
		return opts
	}

	if err := decodeCfgToOpts(cfg, opts); err != nil {
		// todo: logs
	}

	return opts
}

func decodeCfgToOpts(cfg any, opts any) (err error) {
	if reflect.ValueOf(cfg).IsNil() {
		return errors.Errorf("cfg is nil")
	}
	buf, err := sonic.Marshal(cfg)
	if err != nil {
		return errors.WithStack(err)
	}
	if err = sonic.Unmarshal(buf, &opts); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func defaultOpts(name, typ string) any {
	switch typ {
	case MTypeCounter:
		return &prometheus.CounterOpts{Name: name}
	case MTypeGauge:
		return &prometheus.GaugeOpts{Name: name}
	case MTypeHistogram:
		return &prometheus.HistogramOpts{Name: name}
	case MTypeSummary:
		return &prometheus.SummaryOpts{Name: name}
	}
	return nil
}

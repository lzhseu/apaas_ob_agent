package prometheus

import (
	"fmt"
	"runtime/debug"
	"sync"

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

func (c *Collector) Collect(tags map[string]string, value float64) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("goroutine panic: %+v\n%s\n", r, debug.Stack())
			err = errors.Errorf("collect panic: %+v\n", r)
		}
	}()

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
		return errors.Errorf("collect error: invalid metric type: %s", c.Type)
	}
	return nil
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
	var err error
	switch cfg.Type {
	case MTypeCounter:
		collector, err = createCounterCollector(cfg)
	case MTypeGauge:
		collector, err = createGaugeCollector(cfg)
	case MTypeHistogram:
		collector, err = createHistogramCollector(cfg)
	case MTypeSummary:
		collector, err = createSummaryCollector(cfg)
	default:
		return nil, errors.Errorf("[createCollector] invalid metric type: %s", cfg.Type)
	}

	if err != nil {
		return nil, err
	}

	if err = prometheus.Register(collector.prom); err != nil {
		return nil, errors.WithStack(err)
	}

	return collector, nil
}

func createCounterCollector(cfg *conf.PrometheusCfg) (*Collector, error) {
	opts := prometheus.CounterOpts{Name: cfg.Name}
	return &Collector{
		Name: cfg.Name,
		Type: MTypeCounter,
		prom: prometheus.NewCounterVec(opts, cfg.LabelNames),
	}, nil
}

func createGaugeCollector(cfg *conf.PrometheusCfg) (*Collector, error) {
	opts := prometheus.GaugeOpts{Name: cfg.Name}
	return &Collector{
		Name: cfg.Name,
		Type: MTypeGauge,
		prom: prometheus.NewGaugeVec(opts, cfg.LabelNames),
	}, nil
}

func createHistogramCollector(cfg *conf.PrometheusCfg) (*Collector, error) {
	opts := prometheus.HistogramOpts{Name: cfg.Name}

	if len(cfg.Buckets) == 0 {
		opts.Buckets = prometheus.DefBuckets
	}
	if cfg.NativeHistogramBucketFactor != nil {
		opts.NativeHistogramBucketFactor = *cfg.NativeHistogramBucketFactor
	}
	if cfg.NativeHistogramZeroThreshold != nil {
		opts.NativeHistogramZeroThreshold = *cfg.NativeHistogramZeroThreshold
	}
	if cfg.NativeHistogramMaxBucketNumber != nil {
		opts.NativeHistogramMaxBucketNumber = *cfg.NativeHistogramMaxBucketNumber
	}
	if cfg.NativeHistogramMinResetDuration != nil {
		opts.NativeHistogramMinResetDuration = *cfg.NativeHistogramMinResetDuration
	}
	if cfg.NativeHistogramMaxZeroThreshold != nil {
		opts.NativeHistogramMaxZeroThreshold = *cfg.NativeHistogramMaxZeroThreshold
	}
	if cfg.NativeHistogramMaxExemplars != nil {
		opts.NativeHistogramMaxExemplars = *cfg.NativeHistogramMaxExemplars
	}
	if cfg.NativeHistogramExemplarTTL != nil {
		opts.NativeHistogramExemplarTTL = *cfg.NativeHistogramExemplarTTL
	}

	return &Collector{
		Name: cfg.Name,
		Type: MTypeHistogram,
		prom: prometheus.NewHistogramVec(opts, cfg.LabelNames),
	}, nil
}

func createSummaryCollector(cfg *conf.PrometheusCfg) (*Collector, error) {
	opts := prometheus.SummaryOpts{Name: cfg.Name}
	if len(opts.Objectives) == 0 {
		opts.Objectives = cfg.Objectives
	}
	if cfg.MaxAge != nil {
		opts.MaxAge = *cfg.MaxAge
	}
	if cfg.AgeBuckets != nil {
		opts.AgeBuckets = *cfg.AgeBuckets
	}
	if cfg.BufCap != nil {
		opts.BufCap = *cfg.BufCap
	}

	return &Collector{
		Name: cfg.Name,
		Type: MTypeSummary,
		prom: prometheus.NewSummaryVec(opts, cfg.LabelNames),
	}, nil
}

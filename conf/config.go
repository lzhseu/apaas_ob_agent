package conf

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	globalConfig *Config
)

func MustInit() {
	data, err := os.ReadFile("conf/config.yaml")
	if err != nil {
		panic(err)
	}

	globalConfig = &Config{}
	if err = yaml.Unmarshal(data, globalConfig); err != nil {
		panic(err)
	}
}

func GetConfig() *Config {
	return globalConfig
}

type Config struct {
	HttpServerCfg *HttpServerCfg            `yaml:"http_server"` // http 服务配置
	FeishuAppCfg  *FeishuAppCfg             `yaml:"feishu_app"`  // 飞书应用配置
	InnerLogsCfg  *LogsCfg                  `yaml:"inner_logs"`  // agent 自身的日志配置
	PrometheusCfg map[string]*PrometheusCfg `yaml:"prometheus"`  // Prometheus 指标配置, 指标名 -> 配置。如果配置在 yaml 文件中，优先以配置为准，否则使用 prometheus 的默认配置
}

type HttpServerCfg struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type FeishuAppCfg struct {
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
}

type LogsCfg struct {
	LogLevel *string        `yaml:"log_level"`
	Console  *LogConsoleCfg `yaml:"console"`
	File     *LogFileCfg    `yaml:"file"`
	Loki     *LogLokiCfg    `yaml:"loki"`
}

type LogConsoleCfg struct {
	Enable bool `yaml:"enable"`
}

type LogFileCfg struct {
	Enable     bool    `yaml:"enable"`
	FileName   string  `yaml:"filename"`
	MaxAge     *int    `yaml:"max_age"`
	SegDur     *string `yaml:"seg_dur"`      // 按照时间周期切割日志文件，取值：hour、day、week、month
	SegMaxSize *int    `yaml:"seg_max_size"` // 按照文件大小切割日志文件，单位：MB
}

type LogLokiCfg struct {
	Enable bool              `yaml:"enable"`
	Schema string            `yaml:"schema"`
	Host   string            `yaml:"host"`
	Port   string            `yaml:"port"`
	Labels map[string]string `yaml:"labels"`
}

type PrometheusCfg struct {
	Name       string   `yaml:"name"`
	Type       string   `yaml:"type"`
	LabelNames []string `yaml:"label_names"`

	// histogram 指标的配置，每项配置均可选，具体含义见：https://github.com/prometheus/client_golang/blob/main/prometheus/histogram.go#L365
	Buckets                         []float64      `yaml:"buckets" json:"buckets,omitempty"`
	NativeHistogramBucketFactor     *float64       `yaml:"native_histogram_bucket_factor" json:"native_histogram_bucket_factor,omitempty"`
	NativeHistogramZeroThreshold    *float64       `yaml:"native_histogram_zero_threshold" json:"native_histogram_zero_threshold,omitempty"`
	NativeHistogramMaxBucketNumber  *uint32        `yaml:"native_histogram_max_bucket_number" json:"native_histogram_max_bucket_number,omitempty"`
	NativeHistogramMinResetDuration *time.Duration `yaml:"native_histogram_min_reset_duration" json:"native_histogram_min_reset_duration,omitempty"`
	NativeHistogramMaxZeroThreshold *float64       `yaml:"native_histogram_max_zero_threshold" json:"native_histogram_max_zero_threshold,omitempty"`
	NativeHistogramMaxExemplars     *int           `yaml:"native_histogram_max_exemplars" json:"native_histogram_max_exemplars,omitempty"`
	NativeHistogramExemplarTTL      *time.Duration `yaml:"native_histogram_exemplar_ttl" json:"native_histogram_exemplar_ttl,omitempty"`

	// summary 指标的配置，每项配置均可选，具体含义见：https://github.com/prometheus/client_golang/blob/main/prometheus/summary.go#L88
	Objectives map[float64]float64 `yaml:"objectives" json:"objectives,omitempty"`
	MaxAge     *time.Duration      `yaml:"max_age" json:"max_age,omitempty"`
	AgeBuckets *uint32             `yaml:"age_buckets" json:"age_buckets,omitempty"`
	BufCap     *uint32             `yaml:"buf_cap" json:"buf_cap,omitempty"`
}

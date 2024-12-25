package config

import (
	"flag"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	globalConfig   *Config
	configFilePath string
	promSchemaDir  string
	logLevel       string
)

func MustInit() {
	flag.StringVar(&configFilePath, "config-file", "conf/config.yaml", "points to the full path of the configuration file")
	flag.StringVar(&promSchemaDir, "prom-schema-dir", "conf/schema/prometheus", "points to the dir path of the prometheus schema configuration file")
	flag.StringVar(&logLevel, "log-level", "", "log level")
	flag.Parse()

	globalConfig = defaultConfig()

	// 解析 config file
	mustParseFromFile(configFilePath, globalConfig)

	// 解析 prometheus schema
	prometheusCfg := mustParsePrometheusSchema(promSchemaDir)
	globalConfig.PrometheusCfg = prometheusCfg

	// 注入命令行变量
	if logLevel != "" {
		globalConfig.InnerLogsCfg.LogLevel = &logLevel
	}

	// 注入环境变量
	SetEnvVar(EnvKeyAgentServerPort, &globalConfig.HttpServerCfg.Port)
	SetEnvVar(EnvKeyFeishuAPPID, &globalConfig.FeishuAppCfg.AppID)
	SetEnvVar(EnvKeyFeishuAPPSecret, &globalConfig.FeishuAppCfg.AppSecret)
	SetEnvVar(EnvKeyAgentLogLevel, globalConfig.InnerLogsCfg.LogLevel)
	MustSetEnvBoolVar(EnvKeyAgentLogFileEnable, &globalConfig.InnerLogsCfg.File.Enable)
	SetEnvVar(EnvKeyAgentLogFilename, &globalConfig.InnerLogsCfg.File.FileName)
	MustSetEnvBoolVar(EnvKeyAgentLogLokiEnable, &globalConfig.InnerLogsCfg.Loki.Enable)
	SetEnvVar(EnvKeyAgentLogLokiRootURL, &globalConfig.InnerLogsCfg.Loki.RootURL)
}

func mustParseFromFile(filename string, output any) {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(data, output); err != nil {
		panic(err)
	}
}

func mustParsePrometheusSchema(dir string) map[string]*PrometheusCfg {

	// 统一路径格式，后缀带 / ，方便后续处理
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	res := make(map[string]*PrometheusCfg)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()

		if !strings.HasSuffix(filename, ".yaml") && !strings.HasSuffix(filename, "yml") {
			continue
		}

		promSchema := &PrometheusSchema{}
		mustParseFromFile(dir+filename, promSchema)

		prometheusCfg := promSchema.GenCompletePrometheusCfg()
		for k, v := range prometheusCfg {
			res[k] = v
		}
	}

	return res
}

func GetConfig() *Config {
	return globalConfig
}

func defaultConfig() *Config {
	return &Config{
		HttpServerCfg: &HttpServerCfg{},
		FeishuAppCfg:  &FeishuAppCfg{},
		InnerLogsCfg: &LogsCfg{
			Console: &LogConsoleCfg{Enable: false},
			File:    &LogFileCfg{Enable: true},
			Loki:    &LogLokiCfg{Enable: true},
		},
	}
}

type Config struct {
	HttpServerCfg *HttpServerCfg            `yaml:"http_server"` // http 服务配置
	FeishuAppCfg  *FeishuAppCfg             `yaml:"feishu_app"`  // 飞书应用配置
	InnerLogsCfg  *LogsCfg                  `yaml:"inner_logs"`  // agent 自身的日志配置
	PrometheusCfg map[string]*PrometheusCfg `yaml:"-"`           // Prometheus 指标配置, 指标名 -> 配置。从 schema/prometheus 目录下读取 yaml 文件，如果配置在 yaml 文件中，优先以配置为准，否则使用 prometheus 的默认配置
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
	Enable  bool              `yaml:"enable"`
	RootURL string            `yaml:"root_url"`
	Labels  map[string]string `yaml:"labels"`
}

type PrometheusCfg struct {
	Name       string   `yaml:"name"`
	Help       string   `yaml:"help"`
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

// PrometheusSchema 配置文件中定义的 prometheus schema，由 schema 解析出 PrometheusCfg
type PrometheusSchema struct {
	CommonLabelNames []string                  `yaml:"common_label_names"`
	Schema           map[string]*PrometheusCfg `yaml:"schema"`
}

func (p *PrometheusSchema) GenCompletePrometheusCfg() map[string]*PrometheusCfg {
	for _, cfg := range p.Schema {
		cfg.LabelNames = append(p.CommonLabelNames, cfg.LabelNames...)
	}
	return p.Schema
}

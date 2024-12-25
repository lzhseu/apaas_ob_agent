package config

import (
	"os"
	"strconv"
)

const (
	EnvKeyAgentServerPort     = "AGENT_SERVER_PORT"
	EnvKeyAgentLogLevel       = "AGENT_LOG_LEVEL"
	EnvKeyAgentLogFileEnable  = "AGENT_LOG_FILE_ENABLE"
	EnvKeyAgentLogFilename    = "AGENT_LOG_FILE_FILENAME"
	EnvKeyAgentLogLokiEnable  = "AGENT_LOG_LOKI_ENABLE"
	EnvKeyAgentLogLokiRootURL = "AGENT_LOG_LOKI_ROOT_URL"
	EnvKeyFeishuAPPID         = "FEISHU_APP_ID"
	EnvKeyFeishuAPPSecret     = "FEISHU_APP_SECRET"
)

func SetEnvVar(key string, val *string) {
	v, ok := os.LookupEnv(key)
	if ok {
		*val = v
	}
}

func MustSetEnvBoolVar(key string, val *bool) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return
	}
	boolVar, err := strconv.ParseBool(v)
	if err != nil {
		panic(err)
	}
	*val = boolVar
}

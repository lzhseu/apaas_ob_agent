package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/lzhseu/apaas_ob_agent/inner/logs"
	"github.com/lzhseu/apaas_ob_agent/service/alert"
)

func AlertWebhook(w http.ResponseWriter, r *http.Request) {
	// todo: 告警设计
	// 目前先简单配置一个聊天机器人 + 打印日志

	var err error
	defer func() {
		if err != nil {
			logs.Error(fmt.Sprintf("alert webhook handler exec error: %v", err))
		}
	}()

	chatBotID := r.URL.Query().Get("chat_bot_id")
	if chatBotID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(w, `{"message": "chat_bot_id is required"}`)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(w, `{"message": "failed to read request body"}`)
		return
	}
	defer r.Body.Close()

	var data alert.GrafanaAlertData
	if err = sonic.Unmarshal(body, &data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(w, `{"message": "failed to unmarshal request body"}`)
		return
	}

	if err = alert.FeishuAlert(context.Background(), &alert.FeishuAlertParam{
		ChatBotID: chatBotID,
		Data:      &data,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(w, `{"message": "failed to send alert to feishu"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintln(w, `{"message": "success"}`)
}

func PrometheusExporter(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintln(w, `{"message": "pong"}`)
}

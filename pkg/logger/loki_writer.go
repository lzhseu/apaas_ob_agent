package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

// LokiWriter loki http api 详见：详见：https://grafana.com/docs/loki/latest/reference/loki-http-api/#ingest-logs
type LokiWriter struct {
	BaseURL string
	Labels  map[string]string
}

func (l *LokiWriter) Write(p []byte) (n int, err error) {
	req := &LokiPushRequest{
		Streams: []StreamItem{
			{
				Stream: l.Labels,
				Values: [][]string{
					{
						fmt.Sprintf("%v", time.Now().UnixNano()), string(p),
					},
				},
			},
		},
	}

	resp, err := resty.New().
		EnableTrace().
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Post(l.BaseURL + "/loki/api/v1/push")

	if err != nil {
		return 0, errors.WithStack(err)
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent {
		return 0, errors.Errorf("status code: %d, body: %s", resp.StatusCode(), string(resp.Body()))
	}

	return len(p), nil
}

type LokiPushRequest struct {
	Streams []StreamItem `json:"streams"`
}

type StreamItem struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

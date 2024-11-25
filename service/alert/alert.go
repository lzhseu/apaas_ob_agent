package alert

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

func FeishuAlert(ctx context.Context, param *FeishuAlertParam) error {
	validate := validator.New()
	err := validate.Struct(param)
	if err != nil {
		return errors.WithStack(err)
	}

	link := param.Data.Alerts[0].GeneratorURL
	req := &feishuMessageRequest{
		MsgType: "text",
		Content: map[string]string{
			"text": fmt.Sprintf("告警了！！！\n告警链接：%v", link),
		},
	}

	resp, err := resty.New().
		EnableTrace().
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Post(fmt.Sprintf("https://open.larkoffice.com/open-apis/bot/v2/hook/%s", param.ChatBotID))

	if err != nil {
		return errors.WithStack(err)
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.Errorf("status code: %d, body: %s", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}

type feishuMessageRequest struct {
	MsgType string            `json:"msg_type"`
	Content map[string]string `json:"content"`
}

type FeishuAlertParam struct {
	ChatBotID string            `json:"chat_bot_id" validate:"required"`
	Data      *GrafanaAlertData `json:"data" validate:"required"`
}

type GrafanaAlertData struct {
	Receiver string  `json:"receiver"`
	Status   string  `json:"status"`
	Alerts   []Alert `json:"alerts" validate:"min=1"`
}

type Alert struct {
	GeneratorURL string `json:"generatorURL"`
}

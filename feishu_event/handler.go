package feishu_event

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bytedance/sonic"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/lzhseu/apaas_ob_agent/utils"
	"github.com/pkg/errors"
)

// FeishuEventPacket 飞书事件传输协议
type FeishuEventPacket struct {
	Schema string                 `json:"schema"`
	Header *larkevent.EventHeader `json:"header"`
	Event  json.RawMessage        `json:"event"`
}

// NewFeishuEventHandler 创建一个飞书事件处理器
func NewFeishuEventHandler(newBizHandlerFunc func() BizHandlerIface, opts ...FeishuEventOptionFunc) func(context.Context, *larkevent.EventReq) error {
	return func(ctx context.Context, req *larkevent.EventReq) error {

		opt := &FeishuEventOption{}
		for _, fn := range opts {
			fn(opt)
		}

		// 1. 解析数据包
		packet := &FeishuEventPacket{}
		if err := sonic.Unmarshal(req.Body, packet); err != nil {
			return errors.WithStack(err)
		}

		// 2. 校验请求来源
		if opt.needVerifyToken {
			if err := checkVerifyToken(ctx, packet.Header, opt.verifyToken); err != nil {
				return err
			}
		}

		// 3. 解析事件内容
		bizHandler := newBizHandlerFunc()
		if bizHandler == nil {
			return errors.Errorf("biz handler is nil")
		}
		if err := bizHandler.Unmarshal(ctx, packet); err != nil {
			return err
		}

		// 4. 校验事件包格式
		if err := bizHandler.Validate(ctx, packet); err != nil {
			return err
		}

		// 5. 异步处理事件
		utils.Go(
			func() {
				if err := bizHandler.Handle(ctx, packet); err != nil {
					// todo: log
					fmt.Printf("\n[lzh][err] handler err: %v\n", err)
				}
			},
		)

		return nil
	}
}

func checkVerifyToken(ctx context.Context, header *larkevent.EventHeader, token string) error {
	if header == nil || header.Token == "" {
		return errors.Errorf("Illegal Request: header is invalid, header: %v", header)
	}
	if token == "" {
		return errors.Errorf("Illegal Request: verify token you set is empty")
	}

	if header.Token == token {
		return nil
	}

	return errors.Errorf("Illegal Request: AppID(%v) VerificationToken(%v) is not equal to given token(%v), please check", header.AppID, header.Token, token)
}

func WithVerifyToken(token string) FeishuEventOptionFunc {
	return func(opt *FeishuEventOption) {
		opt.needVerifyToken = true
		opt.verifyToken = token
	}
}

type FeishuEventOptionFunc func(*FeishuEventOption)

type FeishuEventOption struct {
	needVerifyToken bool
	verifyToken     string
}

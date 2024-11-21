package feishu_event

import (
	"context"
)

// BizHandlerIface 业务处理器，实现此接口来处理不同的飞书事件，例如：metrics，log，trace 等
type BizHandlerIface interface {
	Unmarshal(ctx context.Context, packet *FeishuEventPacket) error
	Validate(ctx context.Context, packet *FeishuEventPacket) error
	Handle(ctx context.Context, packet *FeishuEventPacket) error
}

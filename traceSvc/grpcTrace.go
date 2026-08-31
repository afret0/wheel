package traceSvc

import (
	"context"

	"github.com/afret0/wheel/tool"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

// GrpcClientOption 返回 gRPC 客户端的链路追踪 DialOption。
// 未开启 TRACE 时返回空 option, 不产生任何额外开销。
func GrpcClientOption() grpc.DialOption {
	if !tool.EnvEnabled("TRACE") {
		return grpc.EmptyDialOption{}
	}

	return grpc.WithStatsHandler(&grpcClientStatsHandler{Handler: otelgrpc.NewClientHandler()})
}

// grpcClientStatsHandler 包装 otelgrpc 的 client handler, 只把标准 gRPC 错误
// 记为 span 失败, 业务自定义错误不算失败。
//
// 背景: otelgrpc 对任何非 nil error 一律 SetStatus(codes.Error), 并且在同一次
// stats.End 处理中就调用了 span.End(), 事后再改状态是无效的。因此只能在事件
// 抵达 otelgrpc 之前拦截。
type grpcClientStatsHandler struct {
	stats.Handler
}

func (h *grpcClientStatsHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	end, ok := rs.(*stats.End)
	if !ok || end.Error == nil || !isBusinessErr(end.Error) {
		h.Handler.HandleRPC(ctx, rs)
		return
	}

	// 业务错误: 保留可观测信息, 但不让 otelgrpc 把 span 标红。
	// 必须在委托之前写属性, 因为委托内部会 End 掉 span。
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		span.SetAttributes(
			attribute.Bool("rpc.business_error", true),
			attribute.String("rpc.business_error.message", status.Convert(end.Error).Message()),
		)
	}

	cp := *end
	cp.Error = nil
	h.Handler.HandleRPC(ctx, &cp)
}

// isBusinessErr 判断是否为业务侧自定义错误。
//
// 业务 handler 直接 return errors.New(...) / frameErr, gRPC 无法识别成 status
// error, 传到客户端统一是 codes.Unknown。而框架层面的失败都有明确的标准码:
// panic -> Internal, 连不上 -> Unavailable, 超时 -> DeadlineExceeded 等,
// 这些仍然按失败处理。
//
// 非 status error(如 context.Canceled)保守起见也按失败处理。
func isBusinessErr(err error) bool {
	s, ok := status.FromError(err)
	if !ok {
		return false
	}

	return s.Code() == codes.Unknown
}

package traceSvc

import (
	"context"

	"github.com/afret0/wheel/tool"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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

// GrpcCallerOptions 返回自动向 RPC metadata 写入 caller 的客户端 DialOption。
//
// 这是 tool.GrpcCtx 之外的第二道防线: 未经 GrpcCtx 构造的 ctx 也能带上 caller。
// 覆盖 unary 与 streaming 两种调用。
func GrpcCallerOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(grpcCallerUnaryClientInterceptor),
		grpc.WithChainStreamInterceptor(grpcCallerStreamClientInterceptor),
	}
}

func grpcCallerUnaryClientInterceptor(
	ctx context.Context,
	method string,
	req, reply any,
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	return invoker(ctxWithCaller(ctx), method, req, reply, cc, opts...)
}

func grpcCallerStreamClientInterceptor(
	ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	return streamer(ctxWithCaller(ctx), desc, cc, method, opts...)
}

// ctxWithCaller 在不破坏已有出站 metadata 的前提下覆盖 caller。
func ctxWithCaller(ctx context.Context) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	} else {
		md = md.Copy()
	}

	tool.InjectCallerMD(md)

	return metadata.NewOutgoingContext(ctx, md)
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

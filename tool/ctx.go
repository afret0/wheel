package tool

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

func OpId(ctx context.Context) string {
	opIdValue := ctx.Value("opId")
	opId, ok := opIdValue.(string)
	if !ok {
		return UUIDWithoutHyphen()
	}
	return opId
}

func GrpcCtx(ctx context.Context) context.Context {
	opId := OpId(ctx)

	//md := metadata.Pairs("opid", opId)
	//
	//if md, ok := metadata.FromIncomingContext(ctx); ok {
	//	if val, exists := md["opid"]; exists && len(val) > 0 {
	//		opId = val[0]
	//	} else {
	//		md["opid"] = []string{opId}
	//		ctx = metadata.NewOutgoingContext(ctx, md)
	//	}
	//}

	spanId := UUIDWithoutHyphen()
	opId = fmt.Sprintf("%s-%s", opId, spanId)
	logrus.Infof("convert opId: %s, caller: %s", opId, CallerInfo(2))

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.Pairs()
	}

	md["opid"] = []string{opId}
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	for k, v := range carrier {
		md[strings.ToLower(k)] = []string{v}
	}

	return metadata.NewOutgoingContext(ctx, md)

	// md, ok := metadata.FromIncomingContext(ctx)
	// if !ok {
	// 	md = metadata.Pairs("opid", opId)
	// 	ctx = metadata.NewOutgoingContext(ctx, md)
	// } else {
	// 	if val, exists := md["opid"]; exists && len(val) > 0 {
	// 		opId = val[0]
	// 	} else {
	// 		md["opid"] = []string{opId}
	// 		//newMd := metadata.Join(md, metadata.Pairs("opid", opId))
	// 		//ctx = metadata.NewOutgoingContext(ctx, newMd)
	// 		ctx = metadata.NewOutgoingContext(ctx, md)
	// 	}
	// }

	// return ctx
}

func OpIdWithoutDefault(ctx context.Context) string {
	opIdValue := ctx.Value("opId")
	opId, ok := opIdValue.(string)
	if !ok {
		return ""
	}
	return opId
}

func NewCtxBK() context.Context {
	return context.WithValue(context.Background(), "opId", strings.ReplaceAll(uuid.New().String(), "-", ""))
}

// RenewCtx 基于 ctx 派生一个不随请求取消的新 context, 用于脱离请求生命周期的后台任务。
// 在丢弃 deadline 与 cancel 的同时保留链路信息: opId 用于日志串联, SpanContext 用于
// 链路串联, 否则后台任务产生的 span 会脱离原链路成为孤立的根 trace。
func RenewCtx(ctx context.Context) context.Context {
	opId := OpId(ctx)
	spanId := UUIDWithoutHyphen()
	opId = fmt.Sprintf("%s-%s", opId, spanId)
	logrus.Infof("convert opId: %s, caller: %s", opId, CallerInfo(2))

	return DetachCtx(ctx, opId)
}

// DetachCtx 返回一个仅保留链路信息的后台 context: 不继承 deadline 与 cancel,
// 但继承 opId 与 SpanContext, 使后台 goroutine 的 span 仍挂在原链路上。
func DetachCtx(ctx context.Context, opId string) context.Context {
	c := context.WithValue(context.Background(), "opId", opId)

	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
		c = trace.ContextWithSpanContext(c, sc)
	}

	return c
}

func ConvertOpId(opId string) string {
	spanId := UUIDWithoutHyphen()
	opId = fmt.Sprintf("%s-%s", opId, spanId)
	logrus.Infof("convert opId: %s, caller: %s", opId, CallerInfo(2))

	return opId
}

// TraceCarrier 把当前链路信息序列化成可随消息体投递的 map, 供消息队列等
// 没有标准协议头的场景使用。无有效 span 时返回 nil。
func TraceCarrier(ctx context.Context) map[string]string {
	if !trace.SpanContextFromContext(ctx).IsValid() {
		return nil
	}

	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	if len(carrier) == 0 {
		return nil
	}

	return carrier
}

// CtxWithTraceCarrier 是 TraceCarrier 的逆操作: 消费侧凭 carrier 还原上游链路,
// 使消费端产生的 span 挂在生产端所属的 trace 上。
func CtxWithTraceCarrier(ctx context.Context, carrier map[string]string) context.Context {
	if len(carrier) == 0 {
		return ctx
	}

	return otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(carrier))
}

func CallerInfo(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	fn := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s:%d %s", file, line, fn.Name())
}

func fn() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return ""
	}

	fn := runtime.FuncForPC(pc)
	return fn.Name()
}

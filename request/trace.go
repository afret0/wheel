package request

import (
	"context"
	"net/http"
	"net/url"

	"github.com/afret0/wheel/tool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "request"

// startClientSpan 为出站 HTTP 调用创建 client span。
// 返回的 context 必须用于后续的 Inject 与请求构造, 否则下游拿到的 parentSpanId
// 会指向调用方的 server span 而非本次调用, 链路层级会失真。
func startClientSpan(ctx context.Context, method, rawURL string) (context.Context, trace.Span) {
	if !tool.EnvEnabled("TRACE") {
		return ctx, nil
	}

	ctx, span := otel.Tracer(tracerName).Start(ctx, method+" "+spanURL(rawURL),
		trace.WithSpanKind(trace.SpanKindClient),
	)
	span.SetAttributes(
		attribute.String("http.request.method", method),
		attribute.String("url.full", rawURL),
	)

	return ctx, span
}

// endClientSpan 记录响应状态并结束 span, span 为 nil 时安全返回
func endClientSpan(span trace.Span, statusCode int, err error) {
	if span == nil {
		return
	}
	defer span.End()

	if statusCode > 0 {
		span.SetAttributes(attribute.Int("http.response.status_code", statusCode))
	}
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}

// spanURL 去除 query 与用户信息, 避免 span 名称基数爆炸及敏感参数泄露
func spanURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	return u.Scheme + "://" + u.Host + u.Path
}

// mergeHeader 把调用方传入的自定义头合并进 hd, 但不允许覆盖链路相关的头
func mergeHeader(hd http.Header, headers []http.Header) {
	if len(headers) == 0 {
		return
	}

	for k, v := range headers[0] {
		if k == "opId" || k == "Traceparent" || k == "Tracestate" {
			continue
		}
		hd[k] = v
	}
}

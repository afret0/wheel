// Package grpcClient 提供统一的 gRPC 客户端建连入口。
//
// 直接使用 grpc.NewClient 需要每处手动拼装 wheel 的标准 DialOption, 一旦遗漏,
// 链路追踪与 caller 透传就会静默失效。统一从本包建连可以避免这类遗漏。
package grpcClient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/afret0/wheel/traceSvc"
)

// New 创建带有 wheel 标准 DialOption 的 gRPC 客户端连接。
//
// opts 追加在默认 option 之后, 因此同类 option(如 transport credentials)
// 可以直接覆盖默认值。
func New(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.NewClient(target, append(DefaultOptions(), opts...)...)
}

// DefaultOptions 返回 wheel 的标准客户端 DialOption。
//
// 默认使用 insecure 传输以适配内网服务间调用; 需要 TLS 时由调用方通过 New 的
// opts 传入 grpc.WithTransportCredentials 覆盖。
func DefaultOptions() []grpc.DialOption {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		traceSvc.GrpcClientOption(),
	}

	return append(opts, traceSvc.GrpcCallerOptions()...)
}

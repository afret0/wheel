package tool

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
		md = metadata.Pairs("opid", opId)
		ctx = metadata.NewOutgoingContext(ctx, md)
	} else {
		if val, exists := md["opid"]; exists && len(val) > 0 {
			opId = val[0]
		} else {
			md["opid"] = []string{opId}
			//newMd := metadata.Join(md, metadata.Pairs("opid", opId))
			//ctx = metadata.NewOutgoingContext(ctx, newMd)
			ctx = metadata.NewOutgoingContext(ctx, md)
		}
	}

	return ctx
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

func RenewCtx(ctx context.Context) context.Context {
	opId := OpId(ctx)
	spanId := UUIDWithoutHyphen()
	opId = fmt.Sprintf("%s-%s", opId, spanId)
	logrus.Infof("convert opId: %s, caller: %s", opId, CallerInfo(2))

	return context.WithValue(context.Background(), "opId", opId)
}

func CallerInfo(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	fn := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s:%d %s", filepath.Base(file), line, fn.Name())
}

func fn() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return ""
	}

	fn := runtime.FuncForPC(pc)
	return fn.Name()
}

package tool

import (
	"context"
	"github.com/afret0/wheel/frame"
)

func Uid(ctx context.Context) string {
	header := frame.Header(ctx)
	uid := header.Get("_uid")
	//if uid == "" {
	//	panic("uid no exist")
	//}

	//uid, ok := ctx.Value("_uid").(string)
	//if !ok {
	//	panic("uid not found in context")
	//}
	//if uid == "" {
	//	panic("uid is empty")
	//}

	return uid
}

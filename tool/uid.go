package tool

import "context"

func Uid(ctx context.Context) string {
	//uid := ctx.Value("_uid")
	//return uid.(string)
	uid, ok := ctx.Value("_uid").(string)
	if !ok {
		panic("uid not found in context")
	}
	if uid == "" {
		panic("uid is empty")
	}

	return uid
}

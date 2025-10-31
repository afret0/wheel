package tool

import (
	"context"
	"fmt"
	"testing"

	"golang.org/x/sync/errgroup"
)

func Test_ctx(t *testing.T) {
	opId := "a"
	c := context.WithValue(NewCtxBK(), "opId", opId)
	fmt.Printf("opId: %s, caller: %s\n", opId, CallerInfo(0))

	RenewCtx(c)

	eg := errgroup.Group{}

	eg.Go(func() error {
		RenewCtx(c)
		return nil
	})

	eg.Go(func() error {
		func() {
			RenewCtx(c)
		}()

		return nil
	})

	eg.Wait()
}

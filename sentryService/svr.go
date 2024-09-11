package sentryService

import (
	"github.com/getsentry/sentry-go"
)

var Env = struct {
	Production  string
	Development string
}{
	Production:  "production",
	Development: "development",
}

func Init(opt sentry.ClientOptions) error {
	opt.TracesSampleRate = 1.0
	opt.EnableTracing = true

	err := sentry.Init(opt)
	return err
}

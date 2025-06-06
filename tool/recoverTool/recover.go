package recoverTool

import (
	"fmt"
	"github.com/afret0/wheel/log"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

type RecoverTool struct {
	service string
	env     string

	emailReceiver []string
	emailSvc      EmailSvc

	limit *Limit

	lg *logrus.Logger
}

type RT = RecoverTool

var rt *RecoverTool

type Option struct {
	Service string `json:"service"`
	Env     string `json:"env"`

	EmailReceiver []string `json:"emailReceiver"`
	EmailSvc      EmailSvc `json:"emailSvc"`
}

func GetRecoverTool(opt *Option) *RecoverTool {
	if rt != nil {
		return rt
	}

	if opt == nil {
		panic("option can not be nil")
	}
	if len(opt.EmailReceiver) != 0 && opt.EmailSvc == nil {
		panic("email service can not be nil when email receiver is set")
	}
	if opt.Service == "" {
		panic("service can not be empty")
	}
	if opt.Env == "" {
		panic("env can not be empty")
	}

	rt = &RecoverTool{
		service: opt.Service,
		env:     opt.Env,

		emailReceiver: opt.EmailReceiver,
		emailSvc:      opt.EmailSvc,

		limit: GetLimit(),

		lg: log.GetLogger(),
	}

	return rt
}

func (r *RecoverTool) Recover() {
	ro := recover()
	if ro == nil {
		return
	}

	r.HandleRecover(ro, string(debug.Stack()))
}

func (r *RecoverTool) HandleRecover(ro any, stack string) {

	//fmt.Println(ro)
	//fmt.Println(stack)

	e, ok := ro.(error)
	if !ok {
		r.lg.Errorf("recover tool ro.(error) failed, ro: %v", ro)
		return
	}

	count := r.limit.Count(e.Error())
	//fmt.Printf("count: %d\n", count)
	if count > 0 {
		r.lg.Infof("recover tool already handled this error, count: %d, error: %s", count, e.Error())
		return
	}

	r.limit.Incr(e.Error())

	htmlContent := formatHtml(r.service, e.Error(), stack)

	err := r.emailSvc.Send(r.emailReceiver, fmt.Sprintf(" [ %s ]  %s [ %s ] ", r.service, e.Error(), r.env), htmlContent)
	if err != nil {
		r.lg.Errorf("err: %s", err)
		return
	}
}

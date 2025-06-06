package recoverTool

import (
	"fmt"
	"github.com/afret0/wheel/request"
	"github.com/afret0/wheel/tool"
	"testing"
)

type EmailSvcMock struct{}

type EReq struct {
	From        string   `json:"from"`
	ToL         []string `json:"toL"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	ContentType string   `json:"contentType"`
}

func (e *EmailSvcMock) Send(toL []string, subject, content string) error {
	url := "http://v2.kekeyuyin.com/email.EmailService/SenEmailToMultiRecipientAsync"

	ret := &struct{}{}
	request.Post(tool.NewCtxBK(), ret, url, &EReq{
		From:        "afreto@kekeyuyin.com",
		ToL:         toL,
		Subject:     subject,
		Body:        content,
		ContentType: "text/html",
	})

	return nil
}

type A struct {
	b string
}

func (a *A) Panic() {
	fmt.Printf("panic: %s\n", a.b)

}

func TestRecover(t *testing.T) {
	defer GetRecoverTool(&Option{
		Service:       "test",
		Env:           tool.GetEnv(),
		EmailReceiver: []string{"kongandmarx@163.com"},
		EmailSvc:      &EmailSvcMock{},
	}).Recover()

	a := &A{}
	a = nil

	a.Panic()
}

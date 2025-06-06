package recoverTool

type EmailSvc interface {
	Send(toL []string, subject, content string) error
}

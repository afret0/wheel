package tm

import "github.com/afret0/wheel/tool/timeTool/timeMock"

var (
	SetOption = timeMock.SetOption
	SetTime   = timeMock.SetTime
	Now       = timeMock.Now
)

type Option = timeMock.Option

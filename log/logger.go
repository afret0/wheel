package log

import (
	"context"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func GetLogger() *logrus.Logger {
	if logger != nil {
		return logger
	}

	logger = logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetReportCaller(true)

	prettyPrint := false

	logger.SetFormatter(&logrus.JSONFormatter{PrettyPrint: prettyPrint, TimestampFormat: "2006-01-02 15:04:05"})
	return logger
}

func CtxLogger(ctx context.Context) *logrus.Entry {
	opIdValue := ctx.Value("opId")
	opId, _ := opIdValue.(string)

	return GetLogger().WithFields(logrus.Fields{"opId": opId})
}

func OpId(ctx context.Context) string {
	opIdValue := ctx.Value("opId")
	opId, _ := opIdValue.(string)
	return opId
}

var middleWareLogger *logrus.Logger

func GetMiddleWareLogger() *logrus.Logger {
	if middleWareLogger != nil {
		return middleWareLogger
	}
	middleWareLogger = logrus.New()
	middleWareLogger.SetLevel(logrus.InfoLevel)
	middleWareLogger.SetReportCaller(false)

	prettyPrint := false

	middleWareLogger.SetFormatter(&logrus.JSONFormatter{PrettyPrint: prettyPrint, TimestampFormat: "2006-01-02 15:04:05"})
	return middleWareLogger
}

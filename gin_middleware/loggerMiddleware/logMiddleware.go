package loggerMiddleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var middleWareLogger *logrus.Logger

func GetMiddleWareLogger() *logrus.Logger {
	if middleWareLogger != nil {
		return middleWareLogger
	}
	middleWareLogger = logrus.New()
	middleWareLogger.SetLevel(logrus.InfoLevel)
	middleWareLogger.SetReportCaller(false)

	middleWareLogger.SetFormatter(&logrus.JSONFormatter{PrettyPrint: false, TimestampFormat: "2006-01-02 15:04:05"})
	return middleWareLogger
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func formatStack(stack string) string {
	lines := strings.Split(stack, "\n")
	var formatted []string
	for i := 0; i < len(lines); i += 2 {
		if i+1 < len(lines) {
			file := strings.TrimSpace(lines[i])
			function := strings.TrimSpace(lines[i+1])
			if strings.HasPrefix(file, "goroutine ") {
				formatted = append(formatted, file)
			} else {
				formatted = append(formatted, fmt.Sprintf("%s\n    at %s", function, file))
			}
		}
	}
	return strings.Join(formatted, "\n")
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		opId := c.GetHeader("opId")
		if opId == "" {
			opId = strings.ReplaceAll(uuid.New().String(), "-", "")
		}

		c.Set("opId", opId)

		startT := time.Now()
		req, _ := c.GetRawData()
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(req))
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		reqUri := c.Request.RequestURI
		token := c.GetHeader("token")

		lg := GetMiddleWareLogger().WithFields(logrus.Fields{
			"uri":   reqUri,
			"token": token,
		})

		defer func() {

			if r := recover(); r != nil {
				stackTrace := formatStack(string(debug.Stack()))
				// 记录panic信息
				lg.WithFields(logrus.Fields{
					"panic": r,
					"stack": stackTrace,
				}).Error("Panic occurred")
				// 设置500状态码
				c.AbortWithStatus(http.StatusInternalServerError)
			}

		}()

		c.Next()
		endT := time.Now()
		latencyT := endT.Sub(startT)
		reqMethod := c.Request.Method
		clientIP := c.ClientIP()
		statusCode := c.Writer.Status()
		uid := c.Request.Header.Get("_uid")

		lg.WithFields(logrus.Fields{
			"reqTime":    startT.Format("2006-01-02 15:04:05"),
			"latencyT":   latencyT.Milliseconds(),
			"method":     reqMethod,
			"clientIP":   clientIP,
			"req":        string(req),
			"res":        blw.body.String(),
			"opId":       opId,
			"uid":        uid,
			"statusCode": statusCode,
		}).Info("请求日志")
	}
}

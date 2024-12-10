package loggerMiddleware

import (
	"bytes"
	"fmt"
	"github.com/getsentry/sentry-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
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

type Option struct {
	Service        string   `json:"service"`
	WhiteList      []string `json:"whiteList"`
	ReportToSentry bool     `json:"reportToSentry"`
	RePanic        bool     `json:"rePanic"`
}

func panicMarshal(occurred any, stackTrace, opId string) string {
	s := fmt.Sprintf("Panic occurred: %s,\nopId: %s,\nstackTrace:\n %s", occurred, opId, stackTrace)
	return s
}

func LoggerMiddleware(opts ...*Option) gin.HandlerFunc {
	opt := new(Option)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}

	return func(c *gin.Context) {
		opId := c.GetHeader("opId")
		if opId == "" {
			opId = strings.ReplaceAll(uuid.New().String(), "-", "")
		}

		c.Set("opId", opId)
		c.Request.Header.Set("opId", opId)

		startT := time.Now()
		req, _ := c.GetRawData()
		c.Request.Body = io.NopCloser(bytes.NewBuffer(req))
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		reqUri := c.Request.RequestURI
		token := c.GetHeader("token")

		clientIP := c.ClientIP()

		lg := GetMiddleWareLogger().WithFields(logrus.Fields{
			"uri":      reqUri,
			"token":    token,
			"req":      string(req),
			"opId":     opId,
			"clientIP": clientIP,
			"reqTime":  startT.Format("2006-01-02 15:04:05"),
		})

		for _, uri := range opt.WhiteList {
			if strings.Contains(reqUri, uri) {
				return
			}
		}

		defer func() {

			if r := recover(); r != nil {
				stackTrace := formatStack(string(debug.Stack()))

				p := panicMarshal(r, stackTrace, opId)

				// 记录panic信息
				lg.WithFields(logrus.Fields{
					"panic": r,
					"stack": stackTrace,
					"opId":  opId,
				}).Error(r)

				if opt.ReportToSentry {
					go sentry.CaptureException(fmt.Errorf("%s", p))
				}
				if opt.RePanic {
					panic(p)
				}

				err := status.Errorf(codes.Internal, "Panic occurred: %#v, stack: %s", r, stackTrace)
				// 设置500状态码
				//c.AbortWithStatus(http.StatusInternalServerError)
				c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "message": err.Error()})
				return
			}

		}()

		c.Next()
		endT := time.Now()
		latencyT := endT.Sub(startT)
		statusCode := c.Writer.Status()
		uid := c.Request.Header.Get("_uid")

		lg.WithFields(logrus.Fields{
			"latencyT":   latencyT.Milliseconds(),
			"res":        blw.body.String(),
			"uid":        uid,
			"statusCode": statusCode,
		}).Info("请求日志")
	}
}

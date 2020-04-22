package gaelog

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	logr "github.com/sirupsen/logrus"
)

const (
	logFile = "/var/log/app.log"

	traceContextHeaderName = "X-Cloud-Trace-Context"

	traceIDContext = "glog-traceID"
)

func Middleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := GetContext(r)
	next(w, r.WithContext(ctx))
}

func GetContext(r *http.Request) context.Context {
	return SetupContext(
		r.Context(),
		r,
	)
}

func SetupContext(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(
		ctx,
		traceIDContext,
		traceID(r),
	)
}

func Debugf(ctx context.Context, data interface{}, format string, v ...interface{}) {
	getLogger(ctx, data).Debugf(format, v...)
}

func Printf(ctx context.Context, data interface{}, format string, v ...interface{}) {
	Infof(ctx, data, format, v...)
}

func Infof(ctx context.Context, data interface{}, format string, v ...interface{}) {
	getLogger(ctx, data).Infof(format, v...)
}

func Warnf(ctx context.Context, data interface{}, format string, v ...interface{}) {
	getLogger(ctx, data).Warnf(format, v...)
}

func Errorf(ctx context.Context, data interface{}, format string, v ...interface{}) {
	getLogger(ctx, data).Errorf(format, v...)
}

func Criticalf(ctx context.Context, data interface{}, format string, v ...interface{}) {
	getLogger(ctx, data).Errorf(format, v...)
}

func Fatalf(ctx context.Context, data interface{}, format string, v ...interface{}) {
	getLogger(ctx, data).Fatalf(format, v...)
}

func traceID(r *http.Request) string {
	return fmt.Sprintf(
		"projects/%s/traces/%s",
		os.Getenv("GOOGLE_CLOUD_PROJECT"),
		strings.Split(r.Header.Get(traceContextHeaderName), "/")[0],
	)
}

func getLogger(ctx context.Context, data interface{}) *logr.Entry {
	return logr.WithContext(ctx).WithFields(logr.Fields{
		"context":                      data,
		"logging.googleapis.com/trace": ctx.Value(traceIDContext),
	})
}

func init() {
	logr.SetLevel(logr.DebugLevel)

	logr.SetFormatter(&logr.JSONFormatter{
		FieldMap: logr.FieldMap{
			logr.FieldKeyTime:  "timestamp",
			logr.FieldKeyLevel: "severity",
			logr.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	})

	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("Failed to open application log file (%s): %s", logFile, err)
	}
	logr.SetOutput(f)
}

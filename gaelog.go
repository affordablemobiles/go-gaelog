package gaelog

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	logr "github.com/sirupsen/logrus"
)

const (
	LogFile = "/var/log/app.log"

	traceContextHeaderName = "X-Cloud-Trace-Context"
)

func traceID(r *http.Request) string {
	return fmt.Sprintf(
		"projects/%s/traces/%s",
		os.Getenv("GOOGLE_CLOUD_PROJECT"),
		strings.Split(r.Header.Get(traceContextHeaderName), "/")[0],
	)
}

func Debugf(r *http.Request, data interface{}, format string, v ...interface{}) {
	getLogger(r, data).Debugf(format, v...)
}

func Printf(r *http.Request, data interface{}, format string, v ...interface{}) {
	Infof(r, data, format, v...)
}

func Infof(r *http.Request, data interface{}, format string, v ...interface{}) {
	getLogger(r, data).Infof(format, v...)
}

func Warnf(r *http.Request, data interface{}, format string, v ...interface{}) {
	getLogger(r, data).Warnf(format, v...)
}

func Errorf(r *http.Request, data interface{}, format string, v ...interface{}) {
	getLogger(r, data).Errorf(format, v...)
}

func Criticalf(r *http.Request, data interface{}, format string, v ...interface{}) {
	getLogger(r, data).Errorf(format, v...)
}

func Fatalf(r *http.Request, data interface{}, format string, v ...interface{}) {
	getLogger(r, data).Fatalf(format, v...)
}

func getLogger(r *http.Request, data interface{}) *logr.Entry {
	return logr.WithContext(r.Context()).WithFields(logr.Fields{
		"context":                      data,
		"logging.googleapis.com/trace": traceID(r),
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

	f, err := os.OpenFile(LogFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("Failed to open application log file (%s): %s", LogFile, err)
	}
	logr.SetOutput(f)
}

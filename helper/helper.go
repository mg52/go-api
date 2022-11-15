package helper

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not found"))
}

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("internal server error"))
}

// NewLogger creates new Logger
func NewLogger(environment string) *logrus.Entry {
	l := logrus.New()
	l.Out = os.Stdout
	l.Formatter = &logrus.JSONFormatter{}
	l.Level = logrus.InfoLevel

	return l.WithFields(logrus.Fields{
		"env": environment,
	})
}

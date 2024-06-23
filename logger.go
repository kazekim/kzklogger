package kzklogger

import (
	"fmt"
	"io"
	"os"
	"sync"

	loggerconstant "github.com/kazekim/kzklogger/constant"
	"github.com/kazekim/kzklogger/formatter"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	Name               string
	ServiceID          string
	ServiceInfo        string
	UserID             string
	RequestID          string
	logrus.FieldLogger // Logger instance
	Data               map[string]interface{}
	isInitial          bool
	mux                *sync.Mutex
}

func (g Logger) SetOutput(w io.Writer) *Logger {
	g.FieldLogger = newDefaultLogger(w).WithFields(g.Data)
	return &g
}

func (g Logger) Printf(format string, args ...interface{}) {
	if !g.isInitial {
		g = *New(g.Name)
	}
	g.FieldLogger.Printf(format, args...)
}

func (g Logger) WithError(err error) *Logger {
	return g.WithField(loggerconstant.FieldError, err)
}

func (g Logger) WithField(key string, value interface{}) *Logger {
	if !g.isInitial {
		g = *New(g.Name)
	}
	f := *g.FieldLogger.WithField(key, value)
	if g.mux == nil {
		g.mux = &sync.Mutex{}
	}
	g.mux.Lock()
	g.Data[key] = value
	g.mux.Unlock()
	g.FieldLogger = &f
	return &g
}

func (g Logger) WithRequestID(value string) *Logger {
	return g.WithField(loggerconstant.FieldRequestID, value)
}
func (g Logger) WithServiceID(value string) *Logger {
	return g.WithField(loggerconstant.FieldServiceID, value)
}

func (g Logger) WithServiceInfo(value string) *Logger {
	return g.WithField(loggerconstant.FieldServiceInfo, value)
}

func (g Logger) WithURL(method string, url string) *Logger {
	return g.WithField(loggerconstant.FieldURL, fmt.Sprintf("%s %s", method, url))
}

func (g Logger) WithUserID(ID string) *Logger {
	return g.WithField(loggerconstant.FieldUserID, ID)
}

func New(name string) *Logger {
	return newLogger(name)
}

func newLogger(name string) *Logger {
	var log = &Logger{
		FieldLogger: newDefaultLogger(os.Stdout),
		Name:        name,
		Data:        map[string]interface{}{},
		isInitial:   true,
		mux:         &sync.Mutex{},
	}
	return log.WithServiceID(name)
}

func newDefaultLogger(output io.Writer) *logrus.Logger {
	var log = &logrus.Logger{
		Out:       output,
		Formatter: &formatter.DefaultLogFormatter{},
		Level:     logrus.DebugLevel,
	}
	return log
}

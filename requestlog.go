package requestlog

import (
	"fmt"
	"net/http"
	"sync"
)

type message struct {
	Category string
	Content  string
}

func (m *message) ToString() string {
	return fmt.Sprintf("%s\t%s", m.Category, m.Content)
}

type logger interface {
	Run()
	GetChan() *chan message
}

var activeRequestLoggers map[string]*RequestLogger
var runningChans map[string]*chan message

func init() {
	activeRequestLoggers = make(map[string]*RequestLogger)
	getLoggerChan("local")
}

//note: *defualtRequestLogger is RequestLogger, defualtRequestLogger is not
type RequestLogger struct {
	ProductName      string
	ThrowWhenTooMany bool
	Chan             *chan message
}

/*
Log the attributes whose names are in keys from HTTP request,
if keys == nil, all the attributes will be logged.
It is better to use it in another go routine like

go Log("", req, nil)

category helps to classify the logs, can be empty
*/
func (this *RequestLogger) Log(category string, req *http.Request, keys *map[string]bool) {

}

func getLoggerChan(loggerType string) *chan message {
	c, ok := runningChans[loggerType]
	if !ok {
		var l logger
		if loggerType == "local" {
			l = &localLogger{make(chan message, 5000)}

		} else {
			l = &remoteLogger{make(chan message, 5000)}
		}
		go l.Run()
		c = l.GetChan()
		runningChans[loggerType] = c

	}
	return c
}

func GetRequestLogger(loggerType string, productName string, throwWhenTooMany bool) *RequestLogger {
	var l sync.Locker
	l.Lock()
	rl, ok := activeRequestLoggers[loggerType+"|"+productName]
	if !ok {
		rl = &RequestLogger{productName, throwWhenTooMany, getLoggerChan(loggerType)}
		activeRequestLoggers[loggerType+"|"+productName] = rl
	}
	l.Unlock()
	return rl
}

//X-Forwarded-For

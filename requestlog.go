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

func (m *message) String() string {
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
	runningChans = make(map[string]*chan message)
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
It is better to use it in another go routine like:
go Log("", req, nil)
category helps to classify the logs, can be empty
*/
func (this *RequestLogger) Log(category string, req *http.Request, canonicalHeaderKeys *map[string]bool, formKeys *map[string]bool) {
	
	if canonicalHeaderKeys == nil {
		canonicalHeaderKeys = make(map[string]bool)
	}
	if formKeys == nil {
		formKeys = make(map[string]bool)
	}
	for k, v := range canonicalHeaderKeys {

	}

}

func GetLocalRequestLogger(productName string, throwWhenTooMany bool, logfunc func(...interface{})) *RequestLogger {
	var l sync.Locker
	l.Lock()
	rl, ok := activeRequestLoggers["local|"+productName]
	if !ok {
		log := &localLogger{make(chan message, 50000), logfunc}
		go log.Run()
		c := log.GetChan()
		rl = &RequestLogger{productName, throwWhenTooMany, c}
		activeRequestLoggers["local|"+productName] = rl
	}
	l.Unlock()
	return rl
}

func GetRemoteRequestLogger(productName string, throwWhenTooMany bool, addr string) *RequestLogger {
	var l sync.Locker
	l.Lock()
	rl, ok := activeRequestLoggers[addr+"|"+productName]
	if !ok {
		log := &remoteLogger{make(chan message, 50000), addr}
		go log.Run()
		c := log.GetChan()
		rl = &RequestLogger{productName, throwWhenTooMany, c}
		activeRequestLoggers[addr+"|"+productName] = rl
	}
	l.Unlock()
	return rl
}

//X-Forwarded-For

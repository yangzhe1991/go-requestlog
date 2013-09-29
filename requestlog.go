package requestlog

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type message struct {
	Product string
	Content string
}

func (m *message) String() string {
	return fmt.Sprintf("%s\t%s", m.Product, m.Content)
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
	ProductName string
	Chan        *chan message
}

/*
Log the attributes whose names are in keys map(their value must be true!) from HTTP request.

If keys == nil, all the attributes will be logged; while if len(keys)==0, nothing will be logged.

If there is no header or form key which want to be logged, the logged value will be "NULL".

If there are more than one values in a key, only the first one will be logged.

"category" helps to classify the logs, can be empty

It is better to use this function in another go routine like:

	go Log("", req, nil, nil)

There are frequently recorded parameters not in header or form will always be logged:

	ip("X-Forwarded-For" first, then RemoteAddr's ip)
	...(to be added)
*/
func (this *RequestLogger) Log(category string, req *http.Request, headerKeys map[string]bool, formKeys map[string]bool) {
	var buffer bytes.Buffer
	//log millisecond rather than nano for compatibility with Youdao's request-log in JAVA.
	buffer.WriteString(string(time.Now().UnixNano()/1000) + "\t" + category)
	if headerKeys == nil {
		for k, vs := range req.Header {
			buffer.WriteString("\t" + escape(k) + "=" + escape(vs[0]))
		}
	} else {
		for k, v := range headerKeys {
			if !v {
				continue
			}
			value, ok := req.Header[k]
			buffer.WriteString("\t" + escape(k) + "=")
			if ok {
				buffer.WriteString(escape(value[0]))
			} else {
				buffer.WriteString("NULL")
			}
		}
	}
	req.ParseForm()
	if formKeys == nil {
		for k, vs := range req.Form {
			buffer.WriteString("\t" + escape(k) + "=" + escape(vs[0]))
		}
	} else {
		for k, v := range formKeys {
			if !v {
				continue
			}
			value, ok := req.Form[k]
			buffer.WriteString("\t" + escape(k) + "=")
			if ok {
				buffer.WriteString(escape(value[0]))
			} else {
				buffer.WriteString("NULL")
			}
		}
	}

	m := message{this.ProductName, buffer.String()}
	*this.Chan <- m
}

func escape(s string) string {
	s = strings.Replace(s, "\t", "\\t", 0)
	s = strings.Replace(s, "\n", "\\n", 0)
	s = strings.Replace(s, "\r", "\\r", 0)
	s = strings.Replace(s, "[", "\\[", 0)
	s = strings.Replace(s, "]", "\\]", 0)
	s = strings.Replace(s, "\\", "\\\\", 0)
	return s
}
func GetLocalRequestLogger(productName string, logfunc func(...interface{})) *RequestLogger {
	var l sync.Locker
	l.Lock()
	rl, ok := activeRequestLoggers["local|"+productName]
	if !ok {
		log := &localLogger{make(chan message, 50000), logfunc}
		go log.Run()
		c := log.GetChan()
		rl = &RequestLogger{productName, c}
		activeRequestLoggers["local|"+productName] = rl
	}
	l.Unlock()
	return rl
}

func GetRemoteRequestLogger(productName string, addr string) *RequestLogger {
	var l sync.Mutex
	l.Lock()
	rl, ok := activeRequestLoggers[addr+"|"+productName]
	if !ok {
		log := &remoteLogger{make(chan message, 50000), addr}
		go log.Run()
		c := log.GetChan()
		rl = &RequestLogger{productName, c}
		activeRequestLoggers[addr+"|"+productName] = rl
	}
	l.Unlock()
	return rl
}

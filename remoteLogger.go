package requestlog

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type remoteLogger struct {
	Chan  chan message
	addrs string
}

func (this *remoteLogger) Run() {
	addrs := strings.Split(this.addrs, ",")
	index := 0
	timeout, _ := time.ParseDuration("500ms")
	for {
		m := <-this.Chan
		conn, err := net.DialTimeout("tcp", addrs[index], timeout)
		if err != nil {
			index = (index + 1) % len(addrs)
			continue
		}
		fmt.Fprintln(conn, m.Product)
		fmt.Fprintln(conn, m.Content)
		conn.Close()
		//fmt.Println(m.Content)
	}
}
func (this *remoteLogger) GetChan() *chan message {
	return &this.Chan
}

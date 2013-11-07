package requestlog

import (
	//"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type remoteLogger struct {
	Chan  chan message
	addrs string
}

func (this *remoteLogger) Run(product string) {
	addrs := strings.Split(this.addrs, ",")
	index := 0
	//fmt.Println("remote logger running")
	timeout, _ := time.ParseDuration("500ms")
	for {
		conn, err := net.DialTimeout("tcp", addrs[index], timeout)
		if err != nil {
			index = (index + 1) % len(addrs)
			fmt.Println("remote logger connection build error!", err, addrs[index])
			continue
		}
		//fmt.Println("remote logger conn build successful!")
		fmt.Fprintln(conn, product)
		counter := 100
		for {
			m := <-this.Chan
			_, err := fmt.Fprintln(conn, m.Content)
			fmt.Println("send one message")
			if err != nil {
				fmt.Println("remote logger send message err!", err)
				break
			}
			if counter--; counter <= 0 {
				break
			}
		}
		//fmt.Println("remote logger conn close")
		conn.Close()
	}
}
func (this *remoteLogger) GetChan() *chan message {
	return &this.Chan
}

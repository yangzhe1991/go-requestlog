package requestlog

type remoteLogger struct {
	Chan  chan message
	addrs string
}

func (this *remoteLogger) Run() {
	for {
		//m := <-this.Chan
	}
}
func (this *remoteLogger) GetChan() *chan message {
	return &this.Chan
}

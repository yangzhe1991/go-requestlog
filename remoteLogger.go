package requestlog

type remoteLogger struct {
	Chan chan message
}

func (this *remoteLogger) Run() {
	for {
		//m := <-this.Chan
	}
}
func (this *remoteLogger) GetChan() *chan message {
	return &this.Chan
}

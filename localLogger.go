package requestlog

type localLogger struct {
	Chan chan message
	log  func(...interface{})
}

func (this *localLogger) Run() {
	for {
		//m := <-this.Chan
	}
}
func (this *localLogger) GetChan() *chan message {
	return &this.Chan
}

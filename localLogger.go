package requestlog

type localLogger struct {
	Chan chan message
	log  func(string)
}

func (this *localLogger) Run() {
	for {
		m := <-this.Chan
		this.log(m.String())
		//fmt.Println(m.Content)
	}
}
func (this *localLogger) GetChan() *chan message {
	return &this.Chan
}

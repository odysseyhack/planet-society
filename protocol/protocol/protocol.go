package protocol

type Protocol struct {
	quit        chan struct{}
	Connections chan Conn
}

func NewProtocol(authorization AuthorizationPlugin) *Protocol {
	return &Protocol{}
}

func (p *Protocol) Loop() {

}

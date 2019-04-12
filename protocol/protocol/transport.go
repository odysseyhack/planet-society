package protocol

type Conn interface {
	Read() (*Message, error)
	Write(*Message) error
	Close() error
}

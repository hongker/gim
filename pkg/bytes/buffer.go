package bytes

type SendBuffer interface {
	Write(b []byte)
	Read() []byte
}

type chanSendBuffer struct {
	queue chan []byte
}

func (buffer chanSendBuffer) Write(b []byte) {
	select {
	case buffer.queue <- b:
	default:
	}
}

func (buffer chanSendBuffer) Read() []byte {
	return <-buffer.queue
}

func NewChanSendBuffer(size int) SendBuffer {
	return &chanSendBuffer{queue: make(chan []byte, size)}
}

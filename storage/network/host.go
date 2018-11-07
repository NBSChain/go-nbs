package network

type StreamHandler func(stream Stream)

type Host interface {
	SetStreamHandler(protocol string, handler StreamHandler)

	NewStream(targetPeerId, protocol string)
}

type NbsHost struct {
}

func (host *NbsHost) SetStreamHandler(protocol string, handler StreamHandler) {

}

func (host *NbsHost) NewStream(targetPeerId, protocol string) {

}

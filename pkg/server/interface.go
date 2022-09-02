package server

type Interface interface {
	Serve(stopCh <-chan struct{}) error
}

type HttpServer struct{}

func (s HttpServer) Serve(stopCh <-chan struct{}) error {
	//TODO implement me
	panic("implement me")
}

type GrpcServer struct{}

func (s GrpcServer) Serve(stopCh <-chan struct{}) error {
	//TODO implement me
	panic("implement me")
}

type SocketServer struct{}

func (s SocketServer) Serve(stopCh <-chan struct{}) error {
	//TODO implement me
	panic("implement me")
}

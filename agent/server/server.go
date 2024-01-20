package server

type Server interface {
	Run() error
}

type server struct {
	opt *Option
}

type Option struct {
	Port int
	Path string
}

func New(opts ...func(opt *Option)) Server {
	opt := &Option{
		Port: 8080,
		Path: "/",
	}
	for _, o := range opts {
		o(opt)
	}

	return &server{
		opt: opt,
	}
}

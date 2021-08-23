package service

type Service interface {
	Init()
	Shutdown()
}

func Init(services []Service) {
	for _, s := range services {
		s.Init()
	}
}

func Shutdown(services []Service) {
	for i := len(services); i >= 0; i--  {
		services[i].Shutdown()
	}
}
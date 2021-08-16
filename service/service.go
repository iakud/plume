package service

type Service interface {
	Init()
	Shutdown()
}
package service

//go:generate mockgen -source=service.go -package=service -destination=service_mock.go

type Filer interface {
	GetPath() string
	ReadWithCheck() (string, error)
}

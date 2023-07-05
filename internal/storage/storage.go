package storage

//go:generate mockgen -source=storage.go -package=storage -destination=storage_mock.go

type Saver interface {
	Save(path, content string) error
}

package store

type Store interface {
	GetTestData() (string, error)
}

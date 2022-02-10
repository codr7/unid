package data

type Source interface {
	Scan(dest ...interface{}) error
}

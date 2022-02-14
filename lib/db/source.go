package db

type Source interface {
	Scan(dst...interface{}) error
}

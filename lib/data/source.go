package data

type Source interface {
	Scan(dst...interface{}) error
}

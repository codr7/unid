package db

import (
	"io"
)

type Rel interface {
	WriteRelSql(out io.Writer) error
	RelParams() []interface{}
}


package data

import (
	"io"
)

type Rel interface {
	WriteRelSql(out io.Writer) error
	RelParams() []interface{}
}


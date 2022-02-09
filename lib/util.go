package unid

import (
	"time"
)

func MinTime() time.Time {
	return time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
}

func MaxTime() time.Time {
	return time.Date(9999, 12, 31, 23, 59, 59, 99999999, time.UTC)
}

package model

import (
	"github.com/dustin/go-humanize"
)

func NewSize(v int64) Size {
	return Size(v)
}

type Size int64

func (this Size) Human() string {
	return humanize.Bytes(uint64(this))
}

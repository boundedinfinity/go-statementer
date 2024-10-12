package model

import "github.com/google/uuid"

type ids struct {
	zero uuid.UUID
}

var Ids = ids{}

func (this ids) IsZero(id uuid.UUID) bool {
	return id.String() == this.zero.String()
}

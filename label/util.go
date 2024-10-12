package label

import (
	"errors"

	"github.com/google/uuid"
)

var ErrLabelManagerErr = errors.New("label manager error")

var (
	LabelParseIdFunc = func(_ int, id string) (uuid.UUID, error) {
		return uuid.Parse(id)
	}
)

// func logErrWrapper[T any](item T, err error) T {
// 	if err != nil {
// 		log.Printf(err.Error())
// 	}

// 	return item
// }

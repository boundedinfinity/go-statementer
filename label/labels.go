package label

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// =====================================================================================
// Simple Label
// =====================================================================================

type SimpleLabel struct {
	Parent      *SimpleLabel
	Children    []*SimpleLabel
	Id          uuid.UUID
	Name        string
	Description string
	Count       int
	Checked     bool
	Selected    bool
}

func (this SimpleLabel) Validate() error {
	if len(this.Name) < 2 {
		return &ErrLabelValidationDetails{
			message: "label.Name must be greater than 2 characters",
			label:   this,
		}
	}

	if len(this.Description) > 0 && len(this.Description) < 2 {
		return &ErrLabelValidationDetails{
			message: "label.Description must be empty or greater than 2 characters",
			label:   this,
		}
	}

	return nil
}

// =====================================================================================
// File Persistence Model
// =====================================================================================

type SimpleLabelPersistenceModelV1 struct {
	Id          uuid.UUID `json:"id" yaml:"id"`
	Name        string    `json:"name" yaml:"name"`
	Description string    `json:"description" yaml:"description"`
}

type SimpleLabelPersistenceModelV2 struct {
	Parent      uuid.UUID `json:"parent" yaml:"parent"`
	Id          uuid.UUID `json:"id" yaml:"id"`
	Name        string    `json:"name" yaml:"name"`
	Description string    `json:"description" yaml:"description"`
}

// =====================================================================================
// Errors
// =====================================================================================

var ErrLabelValidation = errors.New("label validation error")

type ErrLabelValidationDetails struct {
	message string
	label   SimpleLabel
}

func (this ErrLabelValidationDetails) Error() string {
	return fmt.Sprintf("%s : %s : %v", ErrLabelValidation.Error(), this.message, this.label)
}

func (this ErrLabelValidationDetails) Unwrap() error {
	return ErrLabelValidation
}

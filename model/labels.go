package model

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// =====================================================================================
// Simple Label
// =====================================================================================

type SimpleLabel struct {
	Parent      *SimpleLabel `json:"parent" yaml:"parent"`
	Id          uuid.UUID    `json:"id" yaml:"id"`
	Name        string       `json:"name" yaml:"name"`
	Description string       `json:"description" yaml:"description"`
	Count       int          `json:"-" yaml:"-"`
	Checked     bool         `json:"-" yaml:"-"`
	Selected    bool         `json:"-" yaml:"-"`
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
	return fmt.Sprintf("%s : %s : %v", ErrFileDescriptorErr.Error(), this.message, this.label)
}

func (this ErrLabelValidationDetails) Unwrap() error {
	return ErrLabelValidation
}

// =====================================================================================
// Companion
// =====================================================================================

var Labels = labels{}

type labels struct{}

func (this labels) M2P(labels ...*SimpleLabel) []SimpleLabelPersistenceModelV2 {
	var persists []SimpleLabelPersistenceModelV2

	for _, label := range labels {
		persist := SimpleLabelPersistenceModelV2{
			Id:          label.Id,
			Name:        label.Name,
			Description: label.Description,
		}

		if label.Parent != nil {
			persist.Parent = label.Id
		}

		persists = append(persists, persist)
	}

	return persists
}

func (this labels) P2M(persists ...SimpleLabelPersistenceModelV2) []*SimpleLabel {
	var labels []*SimpleLabel

	for _, persist := range persists {
		label := SimpleLabel{
			Id:          persist.Id,
			Name:        persist.Name,
			Description: persist.Description,
		}

		if !Ids.IsZero(persist.Parent) {
			label.Parent = &SimpleLabel{Id: persist.Id}
		}

		labels = append(labels, &label)
	}

	return labels
}

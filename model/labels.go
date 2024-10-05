package model

import (
	"errors"
	"fmt"

	"github.com/boundedinfinity/go-commoner/idiomatic/mapper"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
	"github.com/boundedinfinity/rfc3339date"
	"github.com/google/uuid"
)

// =====================================================================================
// Label
// =====================================================================================

func NewFromLabel(label SimpleLabel) *SimpleLabel {
	return &SimpleLabel{
		Name:        label.Name,
		Description: label.Description,
	}
}

type SimpleLabel struct {
	Id          uuid.UUID `json:"id" yaml:"id"`
	Name        string    `json:"name" yaml:"name"`
	Description string    `json:"description" yaml:"description"`
	Count       int       `json:"-" yaml:"-"`
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

func labelNameFilter(label *SimpleLabel, text string) bool {
	return stringer.Contains(label.Name, text)
}

func labelDescriptionFilter(label *SimpleLabel, text string) bool {
	return stringer.Contains(label.Description, text)
}

// =====================================================================================
// Label Errors
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
// DateLabel
// =====================================================================================

type DateLabel struct {
	*SimpleLabel
	Date rfc3339date.Rfc3339Date `json:"date" yaml:"date"`
}

func DateLabels2Labels(datedLabels []DateLabel) []*SimpleLabel {
	var labels []*SimpleLabel

	for _, datedLabel := range datedLabels {
		labels = append(labels, datedLabel.SimpleLabel)
	}

	return labels
}

// =====================================================================================
// ValueLabel
// =====================================================================================

type ValueLabel struct {
	*SimpleLabel
	Value string `json:"value" yaml:"value"`
}

func ValueLabel2Labels(datedLabels []ValueLabel) []*SimpleLabel {
	var labels []*SimpleLabel

	for _, datedLabel := range datedLabels {
		labels = append(labels, datedLabel.SimpleLabel)
	}

	return labels
}

// =====================================================================================
// Labels
// =====================================================================================

type Labels []*SimpleLabel

func (this Labels) filter(text string, fns ...func(*SimpleLabel, string) bool) []*SimpleLabel {
	var found []*SimpleLabel

	for _, label := range this {
		for _, fn := range fns {
			if fn(label, text) {
				found = append(found, label)
			}
		}
	}

	return found
}

func (this Labels) contains(text string, fns ...func(*SimpleLabel, string) bool) bool {
	var found bool

	for _, label := range this {
		for _, fn := range fns {
			if fn(label, text) {
				found = true
				break
			}
		}
	}

	return found
}

func (this Labels) ByName(text string) []*SimpleLabel {
	return this.filter(text, labelNameFilter)
}

func (this Labels) ByDescription(text string) []*SimpleLabel {
	return this.filter(text, labelDescriptionFilter)
}

func (this Labels) ByTerm(text string) []*SimpleLabel {
	return this.filter(text, labelNameFilter, labelDescriptionFilter)
}

// =====================================================================================
// LabelManager
// =====================================================================================

func NewLabelManager() *LabelManager {
	return &LabelManager{
		byId:   map[string]*SimpleLabel{},
		byName: map[string]*SimpleLabel{},
	}
}

type LabelManager struct {
	byId   map[string]*SimpleLabel
	byName map[string]*SimpleLabel
}

func (this *LabelManager) All() Labels {
	return mapper.Values(this.byId)
}

func (this *LabelManager) New(name, description string) (*SimpleLabel, error) {
	return this.Add(&SimpleLabel{Name: name, Description: description})
}

var ErrLabelManagerErr = errors.New("label manager error")

func (this *LabelManager) ById(id uuid.UUID) *SimpleLabel {
	return this.byId[id.String()]
}

func (this *LabelManager) ByName(id uuid.UUID) *SimpleLabel {
	return this.byId[id.String()]
}

func (this *LabelManager) Add(labels ...*SimpleLabel) (*SimpleLabel, error) {
	for _, label := range labels {
		if current, err := this.add(label); err != nil {
			return current, err
		}
	}

	return nil, nil
}

func (this *LabelManager) add(label *SimpleLabel) (*SimpleLabel, error) {
	if err := label.Validate(); err != nil {
		return nil, NewGenericErrorWrapper(label).WithErrs(ErrLabelManagerErr, err)
	}

	var found *SimpleLabel
	var ok bool

	if !uuidIsZero(label.Id) {
		found, ok = this.byId[label.Id.String()]
		// TODO check descriptions
	}

	if !ok {
		found, ok = this.byName[stringer.Lowercase(label.Name)]
		// TODO check descriptions
	}

	if !ok {
		if uuidIsZero(label.Id) {
			label.Id = uuid.New()
		}

		this.byId[label.Id.String()] = label
		this.byName[stringer.Lowercase(label.Name)] = label

		found = label
	}

	found.Count++
	return found, nil
}

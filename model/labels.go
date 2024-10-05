package model

import (
	"errors"
	"fmt"

	"github.com/boundedinfinity/go-commoner/idiomatic/mapper"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
	"github.com/boundedinfinity/rfc3339date"
	"github.com/google/uuid"
)

// Label

func NewFromLabel(label Label) *Label {
	return &Label{
		Name:        label.Name,
		Description: label.Description,
	}
}

type Label struct {
	Id          uuid.UUID `json:"id" yaml:"id"`
	Name        string    `json:"name" yaml:"name"`
	Description string    `json:"description" yaml:"description"`
	Count       int       `json:"-" yaml:"-"`
}

var ErrLabelValidation = errors.New("label validation error")

type ErrLabelValidationDetails struct {
	message string
	label   Label
}

func (this ErrLabelValidationDetails) Error() string {
	return fmt.Sprintf("%s : %s : %v", ErrFileDescriptorErr.Error(), this.message, this.label)
}

func (this ErrLabelValidationDetails) Unwrap() error {
	return ErrLabelValidation
}

func (this Label) Validate() error {
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

func labelNameFilter(label *Label, text string) bool {
	return stringer.Contains(label.Name, text)
}

func labelDescriptionFilter(label *Label, text string) bool {
	return stringer.Contains(label.Description, text)
}

// DateLabel

type DateLabel struct {
	*Label
	Date rfc3339date.Rfc3339Date `json:"date" yaml:"date"`
}

func DateLabels2Labels(datedLabels []DateLabel) []*Label {
	var labels []*Label

	for _, datedLabel := range datedLabels {
		labels = append(labels, datedLabel.Label)
	}

	return labels
}

// Labels

type Labels []*Label

func (this Labels) filter(text string, fns ...func(*Label, string) bool) []*Label {
	var found []*Label

	for _, label := range this {
		for _, fn := range fns {
			if fn(label, text) {
				found = append(found, label)
			}
		}
	}

	return found
}

func (this Labels) contains(text string, fns ...func(*Label, string) bool) bool {
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

func (this Labels) ByName(text string) []*Label {
	return this.filter(text, labelNameFilter)
}

func (this Labels) ByDescription(text string) []*Label {
	return this.filter(text, labelDescriptionFilter)
}

func (this Labels) ByTerm(text string) []*Label {
	return this.filter(text, labelNameFilter, labelDescriptionFilter)
}

// LabelMap

type LabelMap map[string]*Label

func (this *LabelMap) Reset() {
	*this = map[string]*Label{}
}

func (this LabelMap) Update(labels ...*Label) {
	for _, label := range labels {
		if _, ok := this[label.Name]; !ok {
			this[label.Name] = NewFromLabel(*label)
			this[label.Name].Count = 1
		} else {
			this[label.Name].Count++
		}
	}
}

func (this LabelMap) filter(text string, fns ...func(*Label, string) bool) []*Label {
	var found []*Label

	for _, label := range this {
		for _, fn := range fns {
			if fn(label, text) {
				found = append(found, label)
			}
		}

	}

	return found
}

func (this LabelMap) ByName(text string) []*Label {
	return this.filter(text, labelNameFilter)
}

func (this LabelMap) ByDescription(text string) []*Label {
	return this.filter(text, labelDescriptionFilter)
}

// LabelManager

func NewLabelManager() *LabelManager {
	return &LabelManager{
		byId:   map[string]*Label{},
		byName: map[string]*Label{},
	}
}

type LabelManager struct {
	byId   map[string]*Label
	byName map[string]*Label
}

func (this *LabelManager) All() Labels {
	return mapper.Values(this.byId)
}

func (this *LabelManager) New(name, description string) (*Label, error) {
	return this.Add(&Label{Name: name, Description: description})
}

var ErrLabelManagerErr = errors.New("label manager error")

func (this *LabelManager) ById(id uuid.UUID) *Label {
	return this.byId[id.String()]
}

func (this *LabelManager) ByName(id uuid.UUID) *Label {
	return this.byId[id.String()]
}

func (this *LabelManager) Add(labels ...*Label) (*Label, error) {
	for _, label := range labels {
		if current, err := this.add(label); err != nil {
			return current, err
		}
	}

	return nil, nil
}

func (this *LabelManager) add(label *Label) (*Label, error) {
	if err := label.Validate(); err != nil {
		return nil, NewGenericErrorWrapper(label).WithErrs(ErrLabelManagerErr, err)
	}

	var found *Label
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

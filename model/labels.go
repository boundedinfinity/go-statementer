package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"

	"github.com/google/uuid"
)

// =====================================================================================
// Simple Label
// =====================================================================================

func SimpleLabelCopy(label SimpleLabel) SimpleLabel {
	return SimpleLabel{
		Id:          label.Id,
		Name:        label.Name,
		Description: label.Description,
		Count:       label.Count,
		Checked:     label.Checked,
		Selected:    label.Selected,
	}
}

func SimpleLabelIds(labels []*SimpleLabel) []uuid.UUID {
	var uuids []uuid.UUID

	for _, label := range labels {
		uuids = append(uuids, label.Id)
	}

	return uuids
}

func SimpleLabelsSame(labels []*SimpleLabel, selecteds []uuid.UUID) bool {
	counts := make(map[uuid.UUID]bool, len(labels))

	for _, label := range labels {
		counts[label.Id] = true
	}

	for _, selected := range selecteds {
		if _, ok := counts[selected]; !ok {
			return false
		}
	}

	return true
}

type SimpleLabel struct {
	Id          uuid.UUID `json:"id" yaml:"id"`
	Name        string    `json:"name" yaml:"name"`
	Description string    `json:"description" yaml:"description"`
	Count       int       `json:"-" yaml:"-"`
	Checked     bool      `json:"-" yaml:"-"`
	Selected    bool      `json:"-" yaml:"-"`
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

func labelNameFilter(label SimpleLabel, text string) bool {
	return stringer.Contains(label.Name, text)
}

func labelDescriptionFilter(label SimpleLabel, text string) bool {
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
// LabelManager
// =====================================================================================

func NewLabelManager() *LabelManager {
	return &LabelManager{
		labels: []*SimpleLabel{},
	}
}

type LabelManager struct {
	labels   []*SimpleLabel
	Selected []uuid.UUID
}

func (this *LabelManager) AddSelected(id uuid.UUID) (*SimpleLabel, bool) {
	if !slicer.Exist(id, this.Selected...) {
		this.Selected = append(this.Selected, id)

		if label, ok := this.ById(id); ok {
			label.Selected = true
			return label, ok
		}
	}

	return nil, false
}

func (this *LabelManager) RemoveSelected(id uuid.UUID) (*SimpleLabel, bool) {
	if slicer.Exist(id, this.Selected...) {
		this.Selected = slicer.Filter(func(_ int, current uuid.UUID) bool {
			return id != current
		}, this.Selected...)

		if label, ok := this.ById(id); ok {
			label.Selected = false
			return label, ok
		}
	}
	return nil, false
}

func (this *LabelManager) All() []*SimpleLabel {
	return this.labels
}

var ErrLabelManagerErr = errors.New("label manager error")

func (this *LabelManager) Reset() {
	this.labels = []*SimpleLabel{}
}

func (this *LabelManager) GenerateYear(year int) error {
	var labels []*SimpleLabel

	for month := time.January; month <= time.December; month++ {
		labels = append(labels, &SimpleLabel{
			Name: fmt.Sprintf("%04d.%02d", year, month),
		})
	}

	if err := this.Add(labels...); err != nil {
		return err
	}

	return nil
}

func (this *LabelManager) ById(id uuid.UUID) (*SimpleLabel, bool) {
	for _, label := range this.labels {
		if label.Id == id {
			return label, true
		}
	}

	return nil, false
}

func (this *LabelManager) ByName(name string) (*SimpleLabel, bool) {
	name = stringer.Lowercase(name)

	for _, label := range this.labels {
		if stringer.Lowercase(label.Name) == name {
			return label, true
		}
	}

	return nil, false
}

func (this *LabelManager) Add(labels ...*SimpleLabel) error {
	for _, label := range labels {
		if err := this.add(label); err != nil {
			return err
		}
	}

	return nil
}

func (this *LabelManager) add(label *SimpleLabel) error {
	if label == nil {
		return nil
	}

	if err := label.Validate(); err != nil {
		return NewGenericErrorWrapper(label).WithErrs(ErrLabelManagerErr, err)
	}

	var found *SimpleLabel
	var ok bool

	if !uuidIsZero(label.Id) {
		if found, ok = this.ById(label.Id); ok {
			return nil
		}
	}

	if !ok {
		if found, ok = this.ByName(label.Name); ok {
			return nil
		}
	}
	// TODO check name and descriptions dups

	if !ok {
		found = label

		if uuidIsZero(found.Id) {
			found.Id = uuid.New()
		}

		this.labels = append(this.labels, found)
	} else {
		label = found
	}

	this.labels = slicer.SortFn(func(label *SimpleLabel) string {
		return label.Name
	}, this.labels...)

	return nil
}

func (this *LabelManager) Count(labels ...*SimpleLabel) error {
	for _, label := range labels {
		if err := this.count(label); err != nil {
			return err
		}
	}

	return nil
}

func (this *LabelManager) count(label *SimpleLabel) error {
	if err := label.Validate(); err != nil {
		return NewGenericErrorWrapper(label).WithErrs(ErrLabelManagerErr, err)
	}

	var found *SimpleLabel
	var ok bool

	if uuidIsZero(label.Id) {
		return NewGenericErrorWrapper(label).WithErrs(
			ErrLabelManagerErr,
			fmt.Errorf("label without ID: %+v", label),
		)
	}

	found, ok = this.ById(label.Id)

	if !ok {
		return NewGenericErrorWrapper(label).WithErrs(
			ErrLabelManagerErr,
			fmt.Errorf("no label found with ID: %s", label.Id),
		)
	}

	found.Count++

	return nil
}

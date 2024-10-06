package model

import (
	"errors"
	"fmt"
	"sync"
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
		// byId:   map[string]*SimpleLabel{},
		// byName: map[string]*SimpleLabel{},
	}
}

type LabelManager struct {
	labels []*SimpleLabel
	// byId   map[string]*SimpleLabel
	// byName map[string]*SimpleLabel
	mutex sync.Mutex
}

func (this *LabelManager) All() []SimpleLabel {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return slicer.Map(func(_ int, label *SimpleLabel) SimpleLabel {
		return *label
	}, this.labels...)
}

var ErrLabelManagerErr = errors.New("label manager error")

func (this *LabelManager) GenerateYear(year int) error {
	var labels []SimpleLabel

	for month := time.January; month <= time.December; month++ {
		labels = append(labels, SimpleLabel{
			Name: fmt.Sprintf("%04d.%02d", year, month),
		})
	}

	if _, err := this.Add(false, labels...); err != nil {
		return err
	}

	return nil
}

func (this *LabelManager) ById(id uuid.UUID) (SimpleLabel, bool) {
	for _, label := range this.labels {
		if label.Id == id {
			return *label, true
		}
	}

	return SimpleLabel{}, false
}

func (this *LabelManager) byId(id uuid.UUID) (*SimpleLabel, bool) {
	for _, label := range this.labels {
		if label.Id == id {
			return label, true
		}
	}

	return nil, false
}

func (this *LabelManager) ByName(name string) (SimpleLabel, bool) {
	name = stringer.Lowercase(name)

	for _, label := range this.labels {
		if stringer.Lowercase(label.Name) == name {
			return *label, true
		}
	}

	return SimpleLabel{}, false
}

func (this *LabelManager) byName(name string) (*SimpleLabel, bool) {
	name = stringer.Lowercase(name)

	for _, label := range this.labels {
		if stringer.Lowercase(label.Name) == name {
			return label, true
		}
	}

	return nil, false
}

func (this *LabelManager) Add(addToCount bool, labels ...SimpleLabel) (SimpleLabel, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	for _, label := range labels {
		if current, err := this.add(addToCount, label); err != nil {
			return current, err
		}
	}

	return SimpleLabel{}, nil
}

func (this *LabelManager) add(addToCount bool, label SimpleLabel) (SimpleLabel, error) {
	if err := label.Validate(); err != nil {
		return SimpleLabel{}, NewGenericErrorWrapper(label).WithErrs(ErrLabelManagerErr, err)
	}

	var found *SimpleLabel
	var ok bool

	if !uuidIsZero(label.Id) {
		if found, ok = this.byId(label.Id); ok {
			return *found, nil
		}
	}

	if !ok {
		if found, ok := this.byName(label.Name); ok {
			return *found, nil
		}
	}
	// TODO check name and descriptions dups

	if !ok {
		copy := SimpleLabelCopy(label)
		found = &copy
		if uuidIsZero(found.Id) {
			found.Id = uuid.New()
		}

		this.labels = append(this.labels, found)
	}

	if addToCount {
		found.Count++
	}

	this.labels = slicer.SortFn(func(label *SimpleLabel) string {
		return label.Name
	}, this.labels...)

	return *found, nil
}

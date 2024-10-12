package model

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
	"github.com/google/uuid"
)

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
	if label, ok := this.ById(id); ok {
		this.Selected = append(this.Selected, id)
		label.Checked = true
		return label, ok
	}

	return nil, false
}

func (this *LabelManager) RemoveSelected(id uuid.UUID) (*SimpleLabel, bool) {
	if label, ok := this.ById(id); ok {
		label.Selected = false

		this.Selected = slicer.Filter(func(_ int, current uuid.UUID) bool {
			return id != current
		}, this.Selected...)

		return label, ok
	}
	return nil, false
}

func (this *LabelManager) All() []*SimpleLabel {
	return this.labels
}

func (this *LabelManager) Reset() {
	this.labels = []*SimpleLabel{}
}

func (this *LabelManager) GenerateYearStr(year string) error {
	if yearInt, err := strconv.Atoi(year); err != nil {
		return err
	} else {
		return this.GenerateYear(yearInt)
	}
}

func (this *LabelManager) GenerateYear(year int) error {
	var labels []*SimpleLabel

	yLabel := &SimpleLabel{
		Name: fmt.Sprintf("%04d", year),
	}

	labels = append(labels, yLabel)

	for month := time.January; month <= time.December; month++ {
		labels = append(labels, &SimpleLabel{
			Name:        fmt.Sprintf("%04d.%02d", year, month),
			Description: fmt.Sprintf("%s %d", month.String(), year),
			Parent:      yLabel,
		})
	}

	if err := this.Add(labels...); err != nil {
		return err
	}

	return nil
}

func (this *LabelManager) ByIdStr(id string) (*SimpleLabel, bool) {
	if idP, err := uuid.Parse(id); err != nil {
		// TODO: This should probably panic
		log.Println(err.Error())
		return nil, false
	} else {
		return this.ById(idP)
	}
}

func (this *LabelManager) ById(id uuid.UUID) (*SimpleLabel, bool) {
	for _, label := range this.labels {
		if label.Id == id {
			return label, true
		}
	}

	return nil, false
}

func (this *LabelManager) ResolveUp(id uuid.UUID) ([]*SimpleLabel, bool) {
	var labels []*SimpleLabel

	if label, ok := this.ById(id); ok {
		labels = append(labels, this.resolveUp(label)...)
	}

	return labels, len(labels) > 0
}

func (this *LabelManager) resolveUp(label *SimpleLabel) []*SimpleLabel {
	if label == nil {
		return []*SimpleLabel{}
	}

	return append([]*SimpleLabel{label}, this.resolveUp(label.Parent)...)
}

func (this *LabelManager) ResolveDown(id uuid.UUID) ([]*SimpleLabel, bool) {
	var labels []*SimpleLabel
	current := id

	for {
		label, ok := this.ById(current)

		if !ok {
			break
		}

		if label.Parent != nil && Ids.IsZero(label.Parent.Id) {
			current = label.Parent.Id
		} else {
			break
		}
	}

	return labels, len(labels) > 0
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

	var ok bool

	if !Ids.IsZero(label.Id) {
		if _, ok = this.ById(label.Id); ok {
			return nil
		}
	}

	if !ok {
		if _, ok = this.ByName(label.Name); ok {
			return nil
		}
	}
	//TODO: check name and descriptions dups

	if !ok {
		if Ids.IsZero(label.Id) {
			label.Id = uuid.New()
		}

		this.labels = append(this.labels, label)
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

	if Ids.IsZero(label.Id) {
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

func (this LabelManager) Copy(label SimpleLabel) SimpleLabel {
	return SimpleLabel{
		Id:          label.Id,
		Name:        label.Name,
		Description: label.Description,
		Count:       label.Count,
		Checked:     label.Checked,
		Selected:    label.Selected,
	}
}

func (this LabelManager) Ids(labels []*SimpleLabel) []uuid.UUID {
	return slicer.Map(Label2IdFunc, labels...)
}

func (this LabelManager) IsSame(labels []*SimpleLabel, selecteds []uuid.UUID) bool {
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

func (this *LabelManager) ResolveParents() error {
	for _, label := range this.labels {
		if label.Parent != nil && !Ids.IsZero(label.Parent.Id) {
			if parent, ok := this.ById(label.Parent.Id); ok {
				label.Parent = parent
			}
		}
	}

	return nil
}

// =====================================================================================

var ErrLabelManagerErr = errors.New("label manager error")

var (
	Label2IdFunc = func(_ int, label *SimpleLabel) uuid.UUID {
		return label.Id
	}

	LabelParseIdFunc = func(_ int, id string) (uuid.UUID, error) {
		return uuid.Parse(id)
	}

	labelNameFilter = func(label SimpleLabel, text string) bool {
		return stringer.Contains(label.Name, text)
	}

	labelDescriptionFilter = func(label SimpleLabel, text string) bool {
		return stringer.Contains(label.Description, text)
	}
)

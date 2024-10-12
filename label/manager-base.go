package label

import (
	"fmt"

	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
	"github.com/boundedinfinity/statementer/util"
	"github.com/google/uuid"
)

// =====================================================================================
// LabelManager
// =====================================================================================

func NewLabelManager() *LabelManager {
	return &LabelManager{
		labelList: []*SimpleLabel{},
		labelMap:  make(map[uuid.UUID]*SimpleLabel),
	}
}

type LabelManager struct {
	labelList []*SimpleLabel
	labelMap  map[uuid.UUID]*SimpleLabel
	Selected  []uuid.UUID
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

func (this *LabelManager) Reset() {
	this.labelList = []*SimpleLabel{}
}

func (this *LabelManager) ResolveDown(id uuid.UUID) ([]*SimpleLabel, bool) {
	var labels []*SimpleLabel
	current := id

	for {
		label, ok := this.ById(current)

		if !ok {
			break
		}

		if label.Parent != nil && util.Ids.IsZero(label.Parent.Id) {
			current = label.Parent.Id
		} else {
			break
		}
	}

	return labels, len(labels) > 0
}

func (this *LabelManager) Add(labels ...*SimpleLabel) error {
	for _, label := range labels {
		if err := this.add(label); err != nil {
			return err
		}
	}

	this.labelList = slicer.SortFn(func(label *SimpleLabel) string {
		return label.Name
	}, this.labelList...)

	return nil
}

func (this *LabelManager) add(label *SimpleLabel) error {
	if label == nil {
		return nil
	}

	if err := label.Validate(); err != nil {
		return util.NewGenericErrorWrapper(label).WithErrs(ErrLabelManagerErr, err)
	}

	var ok bool

	if !util.Ids.IsZero(label.Id) {
		if _, ok = this.ById(label.Id); ok {
			return nil
		}
	}

	if !ok {
		if list := this.List(NameEqualsFilter(label.Name)); len(list) > 0 {
			return nil
		}
	}
	//TODO: check name and descriptions dups

	if !ok {
		if util.Ids.IsZero(label.Id) {
			label.Id = uuid.New()
		}

		this.labelList = append(this.labelList, label)
		this.labelMap[label.Id] = label
	}

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
		return util.NewGenericErrorWrapper(label).WithErrs(ErrLabelManagerErr, err)
	}

	var found *SimpleLabel
	var ok bool

	if util.Ids.IsZero(label.Id) {
		return util.NewGenericErrorWrapper(label).WithErrs(
			ErrLabelManagerErr,
			fmt.Errorf("label without ID: %+v", label),
		)
	}

	found, ok = this.ById(label.Id)

	if !ok {
		return util.NewGenericErrorWrapper(label).WithErrs(
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
	return slicer.Map(label2id, labels...)
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

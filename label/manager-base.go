package label

import (
	"fmt"
	"log"

	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
	"github.com/boundedinfinity/statementer/util"
	"github.com/google/uuid"
)

// =====================================================================================
// LabelManager
// =====================================================================================

func NewLabelManager() *LabelManager {
	return &LabelManager{
		labelList: []*LabelViewModel{},
	}
}

type LabelManager struct {
	labelList []*LabelViewModel
}

func (this *LabelManager) Select(state bool, id string) (*LabelViewModel, bool) {
	var label *LabelViewModel
	var found bool

	switch {
	case id == "":
		// nothing to do
	case id == "all":
		for _, label := range this.labelList {
			label.Selected = state
		}

	default:
		if parsed, err := uuid.Parse(id); err == nil {
			if label, found = this.ById(parsed); found {
				label.Selected = state
			}
		} else {
			log.Println(err.Error())
		}
	}

	return label, found
}

func (this *LabelManager) Reset() {
	this.labelList = []*LabelViewModel{}
}

func (this *LabelManager) Add(labels ...*LabelViewModel) error {
	for _, label := range labels {
		if err := this.add(label); err != nil {
			return err
		}
	}

	this.labelList = slicer.SortFn(func(label *LabelViewModel) string {
		return label.Name
	}, this.labelList...)

	return nil
}

func (this *LabelManager) add(label *LabelViewModel) error {
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
	}

	return nil
}

func (this *LabelManager) Count(labels ...*LabelViewModel) error {
	for _, label := range labels {
		if err := this.count(label); err != nil {
			return err
		}
	}

	return nil
}

func (this *LabelManager) count(label *LabelViewModel) error {
	if err := label.Validate(); err != nil {
		return util.NewGenericErrorWrapper(label).WithErrs(ErrLabelManagerErr, err)
	}

	var found *LabelViewModel
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

func (this LabelManager) Copy(label LabelViewModel) LabelViewModel {
	return LabelViewModel{
		Parent:      label.Parent,
		Id:          label.Id,
		Name:        label.Name,
		Description: label.Description,
		Children:    label.Children,
		Count:       label.Count,
		Checked:     label.Checked,
		Selected:    label.Selected,
		Expanded:    label.Expanded,
	}
}

func (this LabelManager) IsSame(labels []*LabelViewModel, selecteds []uuid.UUID) bool {
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

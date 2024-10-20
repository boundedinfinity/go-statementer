package label

import (
	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
	"github.com/boundedinfinity/statementer/util"
	"github.com/google/uuid"
)

func (this *LabelManager) ResolveInit() error {
	for _, label := range this.labelList {
		if label.Parent != nil && !util.Ids.IsZero(label.Parent.Id) {
			if parent, ok := this.ById(label.Parent.Id); ok {
				label.Parent = parent
				parent.Children = append(parent.Children, label)
			}
		}
	}

	for _, label := range this.labelList {
		label.Children = slicer.UniqFn(IdExtract, label.Children...)
	}

	for _, label := range this.labelList {
		label.Children = slicer.SortFn(func(label *LabelViewModel) string {
			return label.Name
		}, label.Children...)
	}

	return nil
}

func (this *LabelManager) ResolveUp(id uuid.UUID) ([]*LabelViewModel, bool) {
	label, _ := this.ById(id)
	labels := this.resolveUp(label)
	return labels, len(labels) > 0
}

func (this *LabelManager) resolveUp(label *LabelViewModel) []*LabelViewModel {
	if label == nil {
		return []*LabelViewModel{}
	}

	return append([]*LabelViewModel{label}, this.resolveUp(label.Parent)...)
}

func (this *LabelManager) ResolveDown(id uuid.UUID) ([]*LabelViewModel, bool) {
	var labels []*LabelViewModel
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

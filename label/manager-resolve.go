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
		label.Children = slicer.UniqFn(label2id, label.Children...)
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

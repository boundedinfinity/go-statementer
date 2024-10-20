package label

import (
	"log"
	"strings"

	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
	"github.com/google/uuid"
)

var (
	IdExtract = func(_ int, label *LabelViewModel) uuid.UUID { return label.Id }
)

// ================================================
// Filters
// ================================================

type LabelFilterFunc func(int, *LabelViewModel) bool

func (this *LabelManager) List(filters ...LabelFilterFunc) []*LabelViewModel {
	results := this.labelList

	for _, filter := range filters {
		results = slicer.Filter(filter, results...)
	}

	results = slicer.SortFn(
		func(label *LabelViewModel) string { return label.Name },
		results...,
	)

	return results
}

func (this *LabelManager) Taxonomy(filters ...LabelFilterFunc) []*LabelViewModel {
	return this.List(append([]LabelFilterFunc{TaxonomyFilter}, filters...)...)
}

var (
	ContainsFilter = func(term string) LabelFilterFunc {
		target := strings.ToLower(term)
		return func(_ int, label *LabelViewModel) bool {
			return stringer.Contains(stringer.Lowercase(label.Name), target) ||
				stringer.Contains(stringer.Lowercase(label.Description), target)
		}
	}

	NameEqualsFilter = func(name string) LabelFilterFunc {
		target := strings.ToLower(name)
		return func(_ int, label *LabelViewModel) bool { return stringer.Lowercase(label.Name) == target }
	}

	IdFilter = func(id uuid.UUID) LabelFilterFunc {
		target := id.String()
		return func(_ int, label *LabelViewModel) bool { return target == label.Id.String() }
	}

	WithoutIdFilter = func(id uuid.UUID) LabelFilterFunc {
		target := id.String()
		return func(_ int, label *LabelViewModel) bool { return target != label.Id.String() }
	}

	WithoutFilter = func(labels ...*LabelViewModel) LabelFilterFunc {
		return func(_ int, label *LabelViewModel) bool {
			for _, filter := range labels {
				if filter != nil && filter.Id == label.Id {
					return false
				}
			}
			return true
		}
	}

	SelectedFilter = func(_ int, label *LabelViewModel) bool { return label.Selected }
	CheckedFilter  = func(_ int, label *LabelViewModel) bool { return label.Checked }
	TaxonomyFilter = func(_ int, label *LabelViewModel) bool { return label.Parent == nil }
)

// ================================================
// Actions
// ================================================

func (this *LabelManager) Each(fns ...func(_ int, label *LabelViewModel)) {
	for _, fn := range fns {
		slicer.Each(fn, this.labelList...)
	}
}

var (
	CheckAction = func(value bool) func(_ int, label *LabelViewModel) {
		return func(_ int, label *LabelViewModel) { label.Checked = value }
	}

	SelectAction = func(value bool) func(_ int, label *LabelViewModel) {
		return func(_ int, label *LabelViewModel) { label.Selected = value }
	}
)

// ================================================
// Find
// ================================================

func (this *LabelManager) ByIdStr(id string) (*LabelViewModel, bool) {
	if idP, err := uuid.Parse(id); err != nil {
		log.Println(err.Error())
		return nil, false
	} else {
		return this.ById(idP)
	}
}

func (this *LabelManager) ById(id uuid.UUID) (*LabelViewModel, bool) {
	label, ok := slicer.FindFn(IdFilter(id), this.labelList...)
	return label, ok
}

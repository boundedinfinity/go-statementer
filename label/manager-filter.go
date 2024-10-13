package label

import (
	"log"

	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
	"github.com/google/uuid"
)

type FilterFunc func(int, *LabelViewModel) bool

func (this *LabelManager) List(filters ...FilterFunc) []*LabelViewModel {
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

func (this *LabelManager) Taxonomy(filters ...FilterFunc) []*LabelViewModel {
	return this.List(append([]FilterFunc{TaxonomyFilter}, filters...)...)
}

func (this *LabelManager) Select(selected bool) {
	slicer.Each(
		func(_ int, label *LabelViewModel) { label.Selected = selected },
		this.labelList...,
	)
}

func (this *LabelManager) Check(checked bool) {
	slicer.Each(
		func(_ int, label *LabelViewModel) { label.Checked = checked },
		this.labelList...,
	)
}

func (this *LabelManager) ByIdStr(id string) (*LabelViewModel, bool) {
	if idP, err := uuid.Parse(id); err != nil {
		log.Println(err.Error())
		return nil, false
	} else {
		return this.ById(idP)
	}
}

func (this *LabelManager) ById(id uuid.UUID) (*LabelViewModel, bool) {
	label, ok := this.labelMap[id]
	return label, ok
}

var (
	ContainsFilter = func(term string) func(_ int, label *LabelViewModel) bool {
		return func(_ int, label *LabelViewModel) bool {
			return nameContainsFilter(*label, term) || descriptionContainsFilter(*label, term)
		}
	}

	NameEqualsFilter = func(name string) func(_ int, label *LabelViewModel) bool {
		return func(_ int, label *LabelViewModel) bool {
			return stringer.Lowercase(label.Name) == stringer.Lowercase(name)
		}
	}

	nameContainsFilter = func(label LabelViewModel, text string) bool {
		return stringer.Contains(label.Name, text)
	}

	descriptionContainsFilter = func(label LabelViewModel, text string) bool {
		return stringer.Contains(label.Description, text)
	}

	// idFilter = func(id uuid.UUID) func(_ int, label *SimpleLabel) bool {
	// 	str1 := id.String()
	// 	return func(_ int, label *SimpleLabel) bool {
	// 		str2 := label.Id.String()
	// 		return str1 == str2
	// 	}
	// }

	TaxonomyFilter = func(_ int, label *LabelViewModel) bool { return label.Parent == nil }

	label2id = func(_ int, label *LabelViewModel) uuid.UUID {
		return label.Id
	}
)

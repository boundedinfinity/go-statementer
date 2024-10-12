package label

import (
	"log"

	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
	"github.com/google/uuid"
)

type FilterFunc func(int, *SimpleLabel) bool

func (this *LabelManager) List(filters ...FilterFunc) []*SimpleLabel {
	results := this.labelList

	for _, filter := range filters {
		results = slicer.Filter(filter, results...)
	}

	return results
}

func (this *LabelManager) Taxonomy(filters ...FilterFunc) []*SimpleLabel {
	return this.List(append([]FilterFunc{taxonomyFilter}, filters...)...)
}

func (this *LabelManager) Select(selected bool) {
	slicer.Each(
		func(_ int, label *SimpleLabel) { label.Selected = selected },
		this.labelList...,
	)
}

func (this *LabelManager) Check(checked bool) {
	slicer.Each(
		func(_ int, label *SimpleLabel) { label.Checked = checked },
		this.labelList...,
	)
}

func (this *LabelManager) ByIdStr(id string) (*SimpleLabel, bool) {
	if idP, err := uuid.Parse(id); err != nil {
		log.Println(err.Error())
		return nil, false
	} else {
		return this.ById(idP)
	}
}

func (this *LabelManager) ById(id uuid.UUID) (*SimpleLabel, bool) {
	label, ok := this.labelMap[id]
	return label, ok
}

var (
	ContainsFilter = func(term string) func(_ int, label *SimpleLabel) bool {
		return func(_ int, label *SimpleLabel) bool {
			return nameContainsFilter(*label, term) || descriptionContainsFilter(*label, term)
		}
	}

	NameEqualsFilter = func(name string) func(_ int, label *SimpleLabel) bool {
		return func(_ int, label *SimpleLabel) bool {
			return stringer.Lowercase(label.Name) == stringer.Lowercase(name)
		}
	}

	nameContainsFilter = func(label SimpleLabel, text string) bool {
		return stringer.Contains(label.Name, text)
	}

	descriptionContainsFilter = func(label SimpleLabel, text string) bool {
		return stringer.Contains(label.Description, text)
	}

	idFilter = func(id uuid.UUID) func(_ int, label *SimpleLabel) bool {
		str1 := id.String()
		return func(_ int, label *SimpleLabel) bool {
			str2 := label.Id.String()
			return str1 == str2
		}
	}

	taxonomyFilter = func(_ int, label *SimpleLabel) bool { return label.Parent == nil }

	label2id = func(_ int, label *SimpleLabel) uuid.UUID {
		return label.Id
	}
)

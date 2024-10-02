package model

import (
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
	"github.com/boundedinfinity/rfc3339date"
)

// Label

func NewFromLabel(label Label) *Label {
	return &Label{
		Name:        label.Name,
		Description: label.Description,
	}
}

type Label struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Count       int    `json:"-" yaml:"-"`
}

func labelNameFilter(label *Label, text string) bool {
	return stringer.Contains(label.Name, text)
}

func labelDescriptionFilter(label *Label, text string) bool {
	return stringer.Contains(label.Description, text)
}

// DateLabel

type DateLabel struct {
	Label
	Date rfc3339date.Rfc3339Date `json:"date" yaml:"date"`
}

func DateLabels2Labels(datedLabels []DateLabel) []*Label {
	var labels []*Label

	for _, datedLabel := range datedLabels {
		labels = append(labels, &datedLabel.Label)
	}

	return labels
}

// Labels

type Labels []*Label

func (this Labels) filter(text string, fns ...func(*Label, string) bool) []*Label {
	var found []*Label

	for _, label := range this {
		for _, fn := range fns {
			if fn(label, text) {
				found = append(found, label)
			}
		}
	}

	return found
}

func (this Labels) contains(text string, fns ...func(*Label, string) bool) bool {
	var found bool

	for _, label := range this {
		for _, fn := range fns {
			if fn(label, text) {
				found = true
				break
			}
		}
	}

	return found
}

func (this Labels) ByName(text string) []*Label {
	return this.filter(text, labelNameFilter)
}

func (this Labels) ByDescription(text string) []*Label {
	return this.filter(text, labelDescriptionFilter)
}

func (this Labels) ByTerm(text string) []*Label {
	return this.filter(text, labelNameFilter, labelDescriptionFilter)
}

// LabelMap

type LabelMap map[string]*Label

func (this *LabelMap) Reset() {
	*this = map[string]*Label{}
}

func (this LabelMap) Update(labels ...*Label) {
	for _, label := range labels {
		if _, ok := this[label.Name]; !ok {
			this[label.Name] = NewFromLabel(*label)
			this[label.Name].Count = 1
		} else {
			this[label.Name].Count++
		}
	}
}

func (this LabelMap) filter(text string, fns ...func(*Label, string) bool) []*Label {
	var found []*Label

	for _, label := range this {
		for _, fn := range fns {
			if fn(label, text) {
				found = append(found, label)
			}
		}

	}

	return found
}

func (this LabelMap) ByName(text string) []*Label {
	return this.filter(text, labelNameFilter)
}

func (this LabelMap) ByDescription(text string) []*Label {
	return this.filter(text, labelDescriptionFilter)
}

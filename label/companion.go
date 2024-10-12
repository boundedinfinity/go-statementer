package label

import "github.com/boundedinfinity/statementer/util"

// =====================================================================================
// Companion
// =====================================================================================

var Labels = labels{}

type labels struct{}

func (this labels) M2P(labels ...*SimpleLabel) []SimpleLabelPersistenceModelV2 {
	var persists []SimpleLabelPersistenceModelV2

	for _, label := range labels {
		persist := SimpleLabelPersistenceModelV2{
			Id:          label.Id,
			Name:        label.Name,
			Description: label.Description,
		}

		if label.Parent != nil {
			persist.Parent = label.Id
		}

		persists = append(persists, persist)
	}

	return persists
}

func (this labels) P2M(persists ...SimpleLabelPersistenceModelV2) []*SimpleLabel {
	var labels []*SimpleLabel

	for _, persist := range persists {
		label := SimpleLabel{
			Id:          persist.Id,
			Name:        persist.Name,
			Description: persist.Description,
		}

		if !util.Ids.IsZero(persist.Parent) {
			label.Parent = &SimpleLabel{Id: persist.Parent}
		}

		labels = append(labels, &label)
	}

	return labels
}

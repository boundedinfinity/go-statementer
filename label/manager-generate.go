package label

import (
	"fmt"
	"strconv"
	"time"
)

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

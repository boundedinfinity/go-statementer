package util

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/boundedinfinity/docsorter/model"
)

func ValidateLineRegex(desc model.LineDescriptor) error {
	if len(desc.Fields) == 0 {
		return nil
	}

	m := make(map[string]bool)

	for _, field := range desc.Fields {
		m[field.Key] = false
	}

	re := regexp.MustCompile(`\(\?P<(?P<named>.*?)\>`)
	foundGroups := re.FindAllStringSubmatch(desc.Pattern, -1)

	for _, foundGroup := range foundGroups {
		if len(foundGroup) == 2 {
			for _, field := range desc.Fields {
				if foundGroup[1] == field.Key {
					m[field.Key] = true
				}
			}
		}
	}

	missing := make([]string, 0)

	for name, ok := range m {
		if !ok {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing [%v] from %v", strings.Join(missing, ","), desc.Pattern)
	}

	return nil
}

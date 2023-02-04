package processors

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/boundedinfinity/go-commoner/stringer"
	"github.com/boundedinfinity/rfc3339date"
	"github.com/oriser/regroup"
)

var (
	usdPattern             = `(?P<usd>[+-]?\$?[\d,]+\.\d{2})`
	accountPattern         = `Account\sNumber:\s*(?P<account>[\d\s]+)`
	openingBalancePatterns = []string{
		`Previous\sBalance\s+` + usdPattern,
		`Beginning\sBalance\s+` + usdPattern,
	}
	closingBalancePatterns = []string{
		`New\sBalance\s+` + usdPattern,
		`Ending\sBalance\s+` + usdPattern,
	}
	openingDatePatterns = []string{
		`(?P<date>\w+\s+\d+,\s+\d+)\s+through`,
		`Opening/Closing Date\s+(?P<date>\d+/\d+/\d+)\s-\s\d+/\d+/\d+`,
	}
	closingDatePatterns = []string{
		`through\s+(?P<date>\w+\s+\d+,\s+\d+)`,
		`Opening/Closing Date\s+\d+/\d+/\d+\s-\s(?P<date>\d+/\d+/\d+)`,
	}

	depositsAndAdditionsBeginPattern = `(?P<depositsBegin>DEPOSITS AND ADDITIONS)`
	depositsAndAdditionsEndPattern   = `(?P<depositsEnd>Total Deposits and Additions)`

	transactionPattern = `(?P<date>\d{2}/\d{2})\s+(?<memo>.*?)\s+` + usdPattern

	chaseCheckDateFormat = "January 02, 2006"
)

func trimLeading0(s string) string { return stringer.TrimLeft(s, "0") }
func removeSpaces(s string) string { return stringer.Remove(s, " ") }
func removeCommas(s string) string { return stringer.Remove(s, ",") }
func removePlus(s string) string   { return stringer.Remove(s, "+") }
func removeDollar(s string) string { return stringer.Remove(s, "$") }

var (
	accountCleanup = []func(string) string{
		strings.TrimSpace,
		removeSpaces,
		trimLeading0,
	}

	dateCleanup = []func(string) string{
		strings.TrimSpace,
		removeSpaces,
	}

	usdCleanup = []func(string) string{
		removeSpaces,
		removeCommas,
		removePlus,
		removeDollar,
	}
)

func matchFn(key string, patterns ...string) func(string) map[string]string {
	res := make([]*regroup.ReGroup, 0)

	for _, pattern := range patterns {
		res = append(res, regroup.MustCompile(pattern))
	}

	return func(line string) map[string]string {
		for _, re := range res {
			m, err := re.Groups(line)

			if err == nil {
				return m
			}
		}

		return make(map[string]string)
	}
}

func containsFn(key string) func(map[string]string) bool {
	return func(m map[string]string) bool {
		_, ok := m[key]
		return ok
	}
}

func cleanFn(key string, fns ...func(string) string) func(map[string]string) {
	return func(m map[string]string) {
		if v, ok := m[key]; ok {
			for _, fn := range fns {
				v = fn(v)
			}

			m[key] = v
		}
	}
}

func extractStringFn(key string, value *string) func(map[string]string) {
	return func(matches map[string]string) {
		if v, ok := matches[key]; ok {
			*value = v
		}
	}
}

func extractFloatFn(key string, value *float32) func(map[string]string) {
	return func(matches map[string]string) {
		var s string
		extractStringFn(key, &s)(matches)

		if s != "" {
			if p, err := strconv.ParseFloat(s, 32); err != nil {
				fmt.Printf("can't parse %v: %v", s, err)
			} else {
				*value = float32(p)
			}
		}
	}
}

func extractDateFn(key string, layout string, value *rfc3339date.Rfc3339Date) func(map[string]string) {
	return func(matches map[string]string) {
		var s string
		extractStringFn(key, &s)(matches)

		if s != "" {
			if parsed, err := time.Parse(layout, s); err != nil {
				fmt.Printf("can't parse %v: %v", s, err)
			} else {
				r := rfc3339date.NewDate(parsed)
				*value = r
			}
		}
	}
}

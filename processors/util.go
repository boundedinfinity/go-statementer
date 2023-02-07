package processors

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/go-commoner/stringer"
	"github.com/boundedinfinity/rfc3339date"
)

var (
	usdPattern         = `(?P<Amount>[+-]?\$?[\d,]+\.\d{2})`
	transactionPattern = `(?P<date>\d{2}/\d{2})\s+(?<memo>.*?)\s+` + usdPattern
	chaseDateFormat1   = "January 02, 2006"
	chaseDateFormat2   = "01/02/2006"

	openingDatePatterns = []string{
		`Opening/Closing Date\s+(?P<date>\d+/\d+/\d+)\s-\s\d+/\d+/\d+`,
	}
	closingDatePatterns = []string{
		`Opening/Closing Date\s+\d+/\d+/\d+\s-\s(?P<date>\d+/\d+/\d+)`,
	}
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
	}

	usdCleanup = []func(string) string{
		removeSpaces,
		removeCommas,
		removePlus,
		removeDollar,
	}
)

func converTransaction(m map[string]string, transaction *model.Transaction, opening, closing rfc3339date.Rfc3339Date) error {
	if err := convertString(m, "Number", &transaction.Number); err != nil {
		if !errors.Is(err, ErrKeyNotFound) {
			return err
		}
	}

	if err := convertString(m, "Memo", &transaction.Memo); err != nil {
		return err
	}

	if err := convertFloat(m, "Amount", &transaction.Amount, usdCleanup...); err != nil {
		return err
	}

	addYear := func() func(string) string {
		return func(s string) string {
			year := opening.Year()

			if opening.Year() != closing.Year() {
				if strings.HasPrefix(s, "01/") || strings.HasPrefix(s, "1/") {
					year = closing.Year()
				}
			}

			s = fmt.Sprintf("%v/%v", s, year)
			return s
		}
	}

	newCleanup := append(dateCleanup, addYear())

	if err := convertDate(m, "Date", chaseDateFormat2, &transaction.Date, newCleanup...); err != nil {
		return err
	}

	return nil
}

var (
	ErrKeyNotFound = errors.New("key not found")
)

func convertString(m map[string]string, key string, v *string, fns ...func(string) string) error {
	if s, ok := m[key]; ok {
		for _, fn := range fns {
			s = fn(s)
		}

		*v = s
		return nil
	}

	return fmt.Errorf("%w %v %v", ErrKeyNotFound, key, m)
}

func convertFloat(m map[string]string, key string, v *float32, fns ...func(string) string) error {
	var s string
	var err error
	var f float64

	if err = convertString(m, key, &s, fns...); err != nil {
		return err
	}

	if f, err = strconv.ParseFloat(s, 32); err != nil {
		return err
	}

	*v = float32(f)
	return nil
}

func convertDate(m map[string]string, key string, layout string, v *rfc3339date.Rfc3339Date, fns ...func(string) string) error {
	var s string
	var err error
	var d time.Time

	if err = convertString(m, key, &s, fns...); err != nil {
		return err
	}

	if d, err = time.Parse(layout, s); err != nil {
		fmt.Printf("can't parse %v: %v", s, err)
	}

	r := rfc3339date.NewDate(d)
	*v = r

	return nil
}

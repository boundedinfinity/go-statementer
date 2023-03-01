package processors

import (
	"fmt"

	"github.com/boundedinfinity/docsorter/model"
)

func (t *ProcessManager) Transform(ocr *model.ProcessContext) error {
	switch ocr.UserConfig.Processor {
	case "chase-checking":
		t.ocr.Checking = model.NewCheckingStatement()
		return t.transformChecking(&t.ocr.Checking)
	case "chase-credit-card":
		t.ocr.CreditCard = model.NewCreditCardStatement()
		t.transformCreditCard(&t.ocr.CreditCard)
	default:
		return fmt.Errorf("error transformer for %v", ocr.UserConfig.Account)
	}

	return nil
}

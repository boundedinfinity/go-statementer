package processors

import (
	"fmt"

	"github.com/boundedinfinity/docsorter/model"
)

func (t *ProcessManager) Transform(pc *model.ProcessContext) error {
	switch pc.UserConfig.Processor {
	case "chase-checking":
		t.pc.Checking = model.NewCheckingStatement()
		return t.transformChecking(&t.pc.Checking)
	case "chase-credit-card":
		t.pc.CreditCard = model.NewCreditCardStatement()
		t.transformCreditCard(&t.pc.CreditCard)
	default:
		return fmt.Errorf("error transformer for %v", pc.UserConfig.Account)
	}

	return nil
}

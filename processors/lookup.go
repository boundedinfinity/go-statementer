package processors

import (
	"fmt"

	"github.com/boundedinfinity/docsorter/model"
)

func (t *ProcessManager) getUserStatementConfig(account string) (model.UserConfigStatement, bool) {
	var config model.UserConfigStatement
	var found bool

	for _, item := range t.userConfig.Statements {
		if item.Account == account {
			config = item
			found = true
			break
		}
	}

	return config, found
}

func (t *ProcessManager) Lookup() (*model.StatementDescriptor, error) {
	var account string

	for _, item := range t.ocr.Data {
		if err := convertString(item.Values, "Account", &account, accountCleanup...); err == nil {
			break
		}
	}

	if account == "" {
		return nil, fmt.Errorf("account not found")
	}

	config, found := t.getUserStatementConfig(account)

	if !found {
		return nil, fmt.Errorf("not user config found for account %v", account)
	}

	t.ocr.UserConfig = config

	var processor *model.StatementDescriptor

	switch config.Processor {
	case "chase-checking":
		processor = t.getChaseChecking()
	default:
		return nil, fmt.Errorf("no processor found for %v/%v", account, config.Processor)
	}

	return processor, nil
}

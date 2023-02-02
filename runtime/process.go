package runtime

import (
	"fmt"

	"github.com/boundedinfinity/docsorter/processors"
)

func (t *Runtime) Process(path string) error {
	account, err := processors.GetAccount(path)

	if err != nil {
		return err
	}

	statement, err := processors.ProcessStatement(path, account)

	if err != nil {
		return err
	}

	// fmt.Printf("               Name: %v\n", sc.Name)
	fmt.Printf("     Account Number: %v\n", account)
	fmt.Printf("     Account Number: %v\n", statement.AccountNumber)
	// fmt.Printf("               Type: %v\n", sc.Processor)
	// fmt.Printf("    Opening Balance: %v\n", openingBalance)
	// fmt.Printf("    Closing Balance: %v\n", closingBalance)

	return nil
}

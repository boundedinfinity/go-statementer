package runtime

import (
	"fmt"

	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/docsorter/processors"
)

func (t *Runtime) Process(path string) error {
	// descriminator, err := processors.Descriminator(path)

	// if err != nil {
	// 	return err
	// }

	descriminator := model.StatementDiscriminator{}

	statement, err := processors.ProcessStatement(path, descriminator)

	if err != nil {
		return err
	}

	// fmt.Printf("               Name: %v\n", sc.Name)
	fmt.Printf("     Account Number: %v\n", descriminator)
	fmt.Printf("     Account Number: %v\n", statement.Account)
	// fmt.Printf("               Type: %v\n", sc.Processor)
	// fmt.Printf("    Opening Balance: %v\n", openingBalance)
	// fmt.Printf("    Closing Balance: %v\n", closingBalance)

	return nil
}

package main

import (
	"fmt"
	"os"

	"github.com/boundedinfinity/docsorter/runtime"
)

// sudo apt install poppler-utils

func main() {
	rt := runtime.New()

	if err := rt.LoadConfig(os.Args[1]); err != nil {
		handleError(err)
		return
	}

	// inputs := []string{
	// "$HOME/Dropbox/statements/bank/chase/6029/20160111-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20180608-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20201109-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20160208-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20180710-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20201208-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20160308-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20180808-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20210111-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20160408-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20180911-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20210208-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20160509-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20181009-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20210308-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20160608-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20181108-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20210408-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20160711-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20181210-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20210510-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20160808-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20190109-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20210608-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20160909-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20190208-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20210709-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20161011-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20190308-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20210809-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20161108-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20190408-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20210909-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20161208-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20190508-statements-6029-.pdf",
	// "$HOME/Dropbox/statements/bank/chase/6029/20211008-statements-6029-.pdf",
	// }

	// for _, input := range inputs {
	// // input := "$HOME/Dropbox/statements/bank/chase/7285/20210123-statements-7285-.pdf"
	// // input := "$HOME/Dropbox/statements/bank/chase/6029/20160111-statements-6029-.pdf"
	// input := "$HOME/Dropbox/statements/bank/chase/6029/20171108-statements-6029-.pdf"

	// fmt.Printf("====================================================================\n")
	// fmt.Printf("               Path: %v\n", input)
	// output, err := rt.OcrSingle(input)

	// if err != nil {
	// 	handleError(err)
	// 	return
	// }

	output := "/tmp/docksorter/20160308-statements-6029-/20160308-statements-6029-.txt"
	// output := "/tmp/docksorter/20160111-statements-6029-/20160111-statements-6029-.txt"
	// err := rt.Process(output)

	err := rt.Process(output)

	if err != nil {
		handleError(err)
		return
	}
	// }
}

func handleError(err error) {
	fmt.Printf("error %v\n", err)
}

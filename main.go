package main

import (
	"fmt"
	"os"

	"github.com/boundedinfinity/docsorter/runtime"
	"github.com/boundedinfinity/docsorter/util"
)

func main() {
	rt := runtime.New()

	if err := rt.LoadUserConfig(os.Args[1]); err != nil {
		handleError(err)
		return
	}

	ocrs, err := rt.LoadFiles()

	if err != nil {
		handleError(err)
		return
	}

	for _, ocr := range ocrs {
		util.PrintSep()
		util.PrintLabeled("Source", ocr.Source)

		if err := rt.OcrSingle(&ocr); err != nil {
			handleError(err)
			return
		}

		if err := rt.Process(&ocr); err != nil {
			handleError(err)
			return
		}
	}
}

func handleError(err error) {
	fmt.Printf("error %v\n", err)
}

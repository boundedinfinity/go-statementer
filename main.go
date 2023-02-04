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

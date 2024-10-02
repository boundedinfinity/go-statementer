package main

import (
	"fmt"
	"os"

	"github.com/boundedinfinity/statementer/runtime"
	"github.com/boundedinfinity/statementer/web"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	rt := runtime.New(logger, os.Args[1])

	if err := rt.LoadConfig(); err != nil {
		handleError(err)
		return
	}

	if err := rt.LoadState(); err != nil {
		handleError(err)
		return
	}

	if err := rt.WalkSource(); err != nil {
		handleError(err)
		return
	}

	// rt.ShowDups()

	w := web.New(rt)

	if err := w.Init(); err != nil {
		handleError(err)
		return
	}

	logger.Fatal(w.Listen())
}

func handleError(err error) {
	fmt.Printf("error %v\n", err)
}

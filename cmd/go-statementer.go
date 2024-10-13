package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/boundedinfinity/statementer/runtime"
	"github.com/boundedinfinity/statementer/web"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	rt := runtime.New(logger)
	var configPath string

	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	if err := rt.LoadConfig(configPath); err != nil {
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

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.Printf("received %s", sig.String())
		err := errors.Join(w.Shutdown(), rt.Shutdown())

		if err != nil {
			log.Println(err.Error())
		}
	}()

	logger.Fatal(w.Listen())
}

func handleError(err error) {
	fmt.Printf("error %v\n", err)
}

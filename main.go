package main

import (
	"fmt"
	"os"

	"github.com/boundedinfinity/docsorter/runtime"
	"github.com/boundedinfinity/docsorter/util"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	rt := runtime.New(logger)

	if err := rt.LoadUserConfig(os.Args[1]); err != nil {
		handleError(err)
		return
	}

	pcs, err := rt.LoadFiles()

	if err != nil {
		handleError(err)
		return
	}

	for _, pc := range pcs {
		logger.Infof(util.PrintSep())
		logger.Infof(util.PrintLabeled("Source", pc.Stage1.Source))

		if err := rt.CalcFiles(rt.UserConfig.WorkPath, "", &pc.Stage1, pc.Stage1); err != nil {
			handleError(err)
			return
		}

		if err := rt.Ocr(&pc.Stage1); err != nil {
			handleError(err)
			return
		}

		if err := rt.Process(&pc); err != nil {
			handleError(err)
			return
		}

		if err := rt.Rename(pc, &pc.Stage2, pc.Stage1); err != nil {
			handleError(err)
			return
		}

		if err := rt.DumpYaml(&pc); err != nil {
			handleError(err)
			return
		}

		if err := rt.DumpCvs(&pc); err != nil {
			handleError(err)
			return
		}

		if err := rt.Output(&pc.Dest, &pc.Stage2); err != nil {
			handleError(err)
			return
		}
	}
}

func handleError(err error) {
	fmt.Printf("error %v\n", err)
}

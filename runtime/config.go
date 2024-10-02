package runtime

import (
	"encoding/json"
	"os"

	"github.com/boundedinfinity/go-commoner/idiomatic/environmenter"
	"github.com/boundedinfinity/go-commoner/idiomatic/pather"
)

func (this *Runtime) LoadConfig() error {
	ev := environmenter.New()

	if err := ev.Process(); err != nil {
		return err
	}

	this.configPath = ev.Substitue(this.configPath)

	data, err := os.ReadFile(this.configPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &this.config); err != nil {
		return err
	}

	this.config.SourceDir = ev.Substitue(this.config.SourceDir)
	this.config.ProcessedDir = ev.Substitue(this.config.ProcessedDir)
	this.statePath = pather.Paths.Join(this.config.ProcessedDir, "state.json")

	return nil
}

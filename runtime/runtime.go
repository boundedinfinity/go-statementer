// Package runtime the runtime
package runtime

import (
	"errors"
	"log"
	"os/exec"
	"runtime"

	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
	"github.com/boundedinfinity/statementer/label"
	"github.com/boundedinfinity/statementer/model"
	"github.com/sirupsen/logrus"
)

func New(logger *logrus.Logger) *Runtime {
	return &Runtime{
		logger: logger,
		Labels: label.NewLabelManager(),
		debug:  true,
	}
}

type Runtime struct {
	Config model.Config
	State  model.StateV1
	Labels *label.LabelManager
	logger *logrus.Logger
	debug  bool
}

func (this *Runtime) OpenRepositoryDir() (string, error) {
	return this.osOpen(this.Config.RepositoryDir)
}

func (this *Runtime) OpenSourceDir() (string, error) {
	return this.osOpen(this.Config.SourceDir)
}

func (this *Runtime) OpenConfigFile() (string, error) {
	return this.osOpen(this.Config.ConfigPath)
}

func (this *Runtime) Debug() bool {
	return this.debug || this.State.Debug || this.Config.Debug
}

func (this *Runtime) osOpen(path string) (string, error) {
	cmds := []string{}

	switch runtime.GOOS {
	case "cygwin", "windows":
		cmds = append(cmds, "cmd", "/c", "start")
	case "linux":
		// cmd = append(cmd, "gnome-open")
		cmds = append(cmds, "xdg-open")
	case "darwin":
		cmds = append(cmds, "open")
	}

	cmds = append(cmds, path)

	if this.Debug() {
		log.Printf("running %s", stringer.Join(" ", cmds...))
	}

	cmd := exec.Command(cmds[0], cmds[1:]...)

	bs, err := cmd.CombinedOutput()
	if err != nil && this.Debug() {
		log.Print(err.Error())
	}

	return string(bs), err
}

func (this *Runtime) Shutdown() error {
	log.Println("shutdown: begin")
	defer log.Println("shutdown: end")
	saveStateErr := this.SaveState()
	return errors.Join(saveStateErr)
}

package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/boundedinfinity/go-commoner/pather"
	"github.com/boundedinfinity/go-commoner/slicer"
)

func (t *Runtime) runDocker(env map[string]string, args []string) (string, error) {
	homeDir := os.Getenv("HOME")
	user := os.Getenv("USER")
	uid := os.Getuid()
	gid := os.Getgid()

	cd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	composeFilePath := pather.Join(cd, "docker/docker-compose.yaml")
	dargs := []string{"--file", composeFilePath, "run"}
	dargs = append(dargs, args...)

	cenv := map[string]string{
		"HOME": homeDir,
		"USER": user,
		"UID":  strconv.FormatInt(int64(uid), 10),
		"GID":  strconv.FormatInt(int64(gid), 10),
	}

	for k, v := range env {
		cenv[k] = v
	}

	denv := []string{}

	for k, v := range cenv {
		denv = append(denv, fmt.Sprintf("%v=%v", k, v))
	}

	docker := exec.Command("docker-compose", dargs...)
	docker.Env = denv

	if t.userConfig.Debug {
		fmt.Printf("%v %v",
			slicer.Join(docker.Env, " "),
			slicer.Join(docker.Args, " "),
		)
	}

	bs, err := docker.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf(
			"while running %v %v\n, failed with: %w\n, output: %v",
			slicer.Join(docker.Env, " "),
			slicer.Join(docker.Args, " "),
			err,
			string(bs),
		)
	}

	return string(bs), nil
}

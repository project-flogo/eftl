package docker

import (
	"io"
	"os"
	"os/exec"
	"time"
)

// StartEFTL starts an EFTL server as a docker image
func StartEFTL() {
	cmd := exec.Command("docker", "start", "eftl")
	stdout, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, stdout)
	err = cmd.Wait()
	if err == nil {
		time.Sleep(60 * time.Second)
		return
	}

	cmd = exec.Command("docker", "run", "--name", "eftl", "-d", "-p", "9191:9191", "-p", "8585:8585", "pointlander/eftl")
	stdout, err = cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, stdout)
	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
	time.Sleep(60 * time.Second)
}

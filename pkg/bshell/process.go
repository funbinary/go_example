package bshell

import (
	"fmt"
	"github.com/funbinary/go_example/pkg/errors"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	envKeyPPid = "GPROC_PPID"
)

type Process struct {
	exec.Cmd
	PPid int
}

// NewProcess creates and returns a new Process.
func NewProcess(path string, args []string, environment ...[]string) *Process {
	env := os.Environ()
	if len(environment) > 0 {
		env = append(env, environment[0]...)
	}
	process := &Process{
		PPid: os.Getpid(),
		Cmd: exec.Cmd{
			Args:       []string{path},
			Path:       path,
			Stdin:      os.Stdin,
			Stdout:     os.Stdout,
			Stderr:     os.Stderr,
			Env:        env,
			ExtraFiles: make([]*os.File, 0),
		},
	}
	process.Dir, _ = os.Getwd()
	if len(args) > 0 {
		// Exclude of current binary path.
		start := 0
		if strings.EqualFold(path, args[0]) {
			start = 1
		}
		process.Args = append(process.Args, args[start:]...)
	}
	return process
}

func (p *Process) Pid() int {
	if p.Process != nil {
		return p.Process.Pid
	}
	return 0
}

func (p *Process) Start() (int, error) {
	if p.Process != nil {
		return p.Pid(), nil
	}
	p.Env = append(p.Env, fmt.Sprintf("%s=%d", envKeyPPid, p.PPid))
	if err := p.Cmd.Start(); err == nil {
		return p.Process.Pid, nil
	} else {
		return 0, err
	}
}

func (p *Process) Run() error {
	if _, err := p.Start(); err == nil {
		return p.Wait()
	} else {
		return err
	}
}

func (p *Process) Kill() (err error) {
	err = p.Process.Kill()
	if err != nil {
		err = errors.Wrapf(err, `kill process failed for pid "%d"`, p.Process.Pid)
		return err
	}
	if runtime.GOOS != "windows" {
		if err = p.Process.Release(); err != nil {
			fmt.Printf(`%+v`, err)
		}
	}
	// It ignores this error, just log it.
	_, err = p.Process.Wait()
	fmt.Printf(`%+v`, err)
	return nil
}

func (p *Process) Signal(sig os.Signal) error {
	return p.Process.Signal(sig)
}

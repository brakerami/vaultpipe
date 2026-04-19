package process

import (
	"errors"
	"os"
	"os/exec"
)

// Runner holds the configuration for executing a subprocess with injected env.
type Runner struct {
	Env  []string
	Args []string
}

// NewRunner creates a Runner with the given environment and command args.
func NewRunner(env []string, args []string) (*Runner, error) {
	if len(args) == 0 {
		return nil, errors.New("process: no command specified")
	}
	return &Runner{Env: env, Args: args}, nil
}

// Run executes the command, replacing the current process via exec.Command.
// Stdout and Stderr are inherited from the parent process.
func (r *Runner) Run() error {
	path, err := exec.LookPath(r.Args[0])
	if err != nil {
		return err
	}

	cmd := exec.Command(path, r.Args[1:]...)
	cmd.Env = r.Env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Packge pubsubemulator provides utilities for Google Pubsub Emulator.
package pubsubemulator

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/shirou/gopsutil/v3/process"
)

// Controller describes a controller for pubsub emulator.
type Controller struct {
	host string
}

// Host returns emulator host.
func (e *Controller) Host() string {
	return e.host
}

// Stop stops emulator processes.
func (e *Controller) Stop(ctx context.Context) error {
	processes, err := process.ProcessesWithContext(ctx)
	if err != nil {
		return err
	}
	for _, p := range processes {
		cmdline, err := p.CmdlineWithContext(ctx)
		if err != nil {
			return err
		}
		if strings.Contains(cmdline, "java") && strings.Contains(cmdline, "pubsub") && strings.Contains(cmdline, "emulator") {
			if err := syscall.Kill(int(p.Pid), syscall.SIGTERM); err != nil {
				return err
			}
		}
	}
	return nil
}

// New returns a new Controller.
func New(ctx context.Context, projectID string) (*Controller, error) {
	cmd := exec.CommandContext(ctx, "gcloud", "beta", "emulators", "pubsub", "start", "--project", projectID)
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	out, err := exec.CommandContext(ctx, "gcloud", "beta", "emulators", "pubsub", "env-init").Output()
	if err != nil {
		return nil, err
	}
	return &Controller{
		host: strings.Split(strings.Split(strings.TrimSpace(string(out)), " ")[1], "=")[1],
	}, nil
}

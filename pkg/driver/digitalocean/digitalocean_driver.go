// SPDX-FileCopyrightText: Copyright The Lima Authors
// SPDX-License-Identifier: Apache-2.0

package digitalocean

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os/exec"
	"path/filepath"

	"github.com/digitalocean/godo"
	"github.com/sirupsen/logrus"

	"github.com/lima-vm/lima/v2/pkg/driver"
	"github.com/lima-vm/lima/v2/pkg/executil"
	"github.com/lima-vm/lima/v2/pkg/limatype"
	"github.com/lima-vm/lima/v2/pkg/ptr"
)

type LimaDigitalOceanDriver struct {
	Instance     *limatype.Instance
	SSHLocalPort int

	qCmd    *exec.Cmd
	qWaitCh chan error

	client *godo.Client
}

var _ driver.Driver = (*LimaDigitalOceanDriver)(nil)

func New() *LimaDigitalOceanDriver {
	return &LimaDigitalOceanDriver{}
}

func (l *LimaDigitalOceanDriver) Configure(inst *limatype.Instance) *driver.ConfiguredDriver {
	l.Instance = inst
	l.SSHLocalPort = inst.SSHLocalPort

	return &driver.ConfiguredDriver{
		Driver: l,
	}
}

func (l *LimaDigitalOceanDriver) Validate(ctx context.Context) error {
	return validateConfig(ctx, l.Instance.Config)
}

func validateConfig(_ context.Context, cfg *limatype.LimaYAML) error {
	if cfg == nil {
		return errors.New("configuration is nil")
	}
	if *cfg.MountType != limatype.REVSSHFS {
		return fmt.Errorf("field `mountType` must be %q for %s driver, got %q",
			limatype.REVSSHFS, "godo", *cfg.MountType)
	}
	return nil
}

func (l *LimaDigitalOceanDriver) FillConfig(ctx context.Context, cfg *limatype.LimaYAML, _ string) error {
	if cfg.VMType == nil {
		cfg.VMType = ptr.Of("digitalocean")
	}
	if cfg.MountType == nil {
		cfg.MountType = ptr.Of(limatype.REVSSHFS)
	}
	return validateConfig(ctx, cfg)
}

func (l *LimaDigitalOceanDriver) BootScripts() (map[string][]byte, error) {
	return nil, nil
}

func (l *LimaDigitalOceanDriver) CreateDisk(_ context.Context) error {
	return nil
}

func (l *LimaDigitalOceanDriver) Start(_ context.Context) (chan error, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		if l.qCmd == nil {
			cancel()
		}
	}()

	qCfg := Config{
		Name:         l.Instance.Name,
		InstanceDir:  l.Instance.Dir,
		LimaYAML:     l.Instance.Config,
		SSHLocalPort: l.SSHLocalPort,
		SSHAddress:   l.Instance.SSHAddress,
	}

	var qArgsFinal []string
	qCmd := exec.CommandContext(ctx, "godo", qArgsFinal...)
	qCmd.SysProcAttr = executil.BackgroundSysProcAttr
	qStdout, err := qCmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	go logPipeRoutine(qStdout, "godo[stdout]")
	qStderr, err := qCmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	go logPipeRoutine(qStderr, "godo[stderr]")

	logrus.Infof("Starting QEMU (hint: to watch the boot progress, see %q)", filepath.Join(qCfg.InstanceDir, "serial*.log"))
	if err := qCmd.Start(); err != nil {
		return nil, err
	}

	l.qWaitCh = make(chan error, 1)

	return l.qWaitCh, nil
}

func (l *LimaDigitalOceanDriver) Stop(_ context.Context) error {
	return errUnimplemented
}

func (l *LimaDigitalOceanDriver) GuestAgentConn(_ context.Context) (net.Conn, string, error) {
	return nil, "", nil
}

func logPipeRoutine(r io.Reader, header string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		logrus.Debugf("%s: %s", header, line)
	}
}

func (l *LimaDigitalOceanDriver) Info() driver.Info {
	var info driver.Info
	info.Name = "digitalocean"
	if l.Instance != nil && l.Instance.Dir != "" {
		info.InstanceDir = l.Instance.Dir
	}

	info.Features = driver.DriverFeatures{
		DynamicSSHAddress:    false,
		SkipSocketForwarding: false,
		CanRunGUI:            false,
	}
	return info
}

func (l *LimaDigitalOceanDriver) SSHAddress(_ context.Context) (string, error) {
	return "127.0.0.1", nil
}

func (l *LimaDigitalOceanDriver) InspectStatus(_ context.Context, _ *limatype.Instance) string {
	return ""
}

func (l *LimaDigitalOceanDriver) RunGUI() error {
	return nil
}

func (l *LimaDigitalOceanDriver) ForwardGuestAgent() bool {
	// if driver is not providing, use host agent
	return true
}

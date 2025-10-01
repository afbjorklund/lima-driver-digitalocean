// SPDX-FileCopyrightText: Copyright The Lima Authors
// SPDX-License-Identifier: Apache-2.0

package digitalocean

import (
	"context"
)

func (l *LimaDigitalOceanDriver) Create(_ context.Context) error {
	return nil
}

func (l *LimaDigitalOceanDriver) Delete(_ context.Context) error {
	return nil
}

func (l *LimaDigitalOceanDriver) Register(_ context.Context) error {
	return nil
}

func (l *LimaDigitalOceanDriver) Unregister(_ context.Context) error {
	return nil
}

func (l *LimaDigitalOceanDriver) ChangeDisplayPassword(_ context.Context, _ string) error {
	return nil
}

func (l *LimaDigitalOceanDriver) DisplayConnection(_ context.Context) (string, error) {
	return "", nil
}

func (l *LimaDigitalOceanDriver) CreateSnapshot(_ context.Context, _ string) error {
	return errUnimplemented
}

func (l *LimaDigitalOceanDriver) ApplySnapshot(_ context.Context, _ string) error {
	return errUnimplemented
}

func (l *LimaDigitalOceanDriver) DeleteSnapshot(_ context.Context, _ string) error {
	return errUnimplemented
}

func (l *LimaDigitalOceanDriver) ListSnapshots(_ context.Context) (string, error) {
	return "", errUnimplemented
}

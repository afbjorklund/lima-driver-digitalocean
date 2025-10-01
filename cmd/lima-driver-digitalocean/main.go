// SPDX-FileCopyrightText: Copyright The Lima Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"

	"github.com/lima-vm/lima/v2/pkg/driver/external/server"
	"github.com/afbjorklund/lima-driver-digitalocean/pkg/driver/digitalocean"
)

// To be used as an external driver for Lima.
func main() {
	server.Serve(context.Background(), digitalocean.New())
}

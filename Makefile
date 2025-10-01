
GO = go

lima-driver-digitalocean: cmd/lima-driver-digitalocean pkg/driver/digitalocean go.mod
	$(GO) build ./cmd/lima-driver-digitalocean

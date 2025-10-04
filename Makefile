
GO = go

PREFIX ?= /usr/local

BIN = lima-driver-digitalocean
CMD = cmd/lima-driver-digitalocean
PKG = pkg/driver/digitalocean
all: $(BIN)

$(BIN): $(CMD) $(PKG) go.mod
	$(GO) build ./$(CMD)

.PHONY: install
install: $(BIN)
	$(INSTALL) -D -m 755 $@ $(DESTDIR)$(PREFIX)/libexec/lima/$(BIN)

.PHONY: lint
lint:
	golangci-lint run $(CMD) $(PKG)

.PHONY: clean
clean:
	$(RM) $(BIN)

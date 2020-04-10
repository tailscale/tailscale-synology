TAILSCALE_VERSION := 0.97-45
TAILSCALE_ARCH=amd64
TAILSCALE_BASE_URL := https://pkgs.tailscale.com/stable
TAILSCALE_PKG := tailscale_$(TAILSCALE_VERSION)_$(TAILSCALE_ARCH).tgz
TAILSCALE_DIR := tailscale_$(TAILSCALE_VERSION)_$(TAILSCALE_ARCH)

all: $(TAILSCALE_DIR)

$(TAILSCALE_DIR): $(TAILSCALE_PKG)
	tar -xzf $(TAILSCALE_PKG)

$(TAILSCALE_PKG):
	wget $(TAILSCALE_BASE_URL)/$(TAILSCALE_PKG)

install: all
	mkdir -p $(DESTDIR)/bin
	install $(TAILSCALE_DIR)/tailscale $(DESTDIR)/bin/
	install $(TAILSCALE_DIR)/tailscaled $(DESTDIR)/bin/

clean:
	rm -rf $(TAILSCALE_DIR) $(TAILSCALE_PKG)

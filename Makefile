TAILSCALE_VERSION := 0.97-45

.PHONY: tailscale-% clean purge

all: tailscale-amd64

tailscale-%:
	@./build-package.sh ${TAILSCALE_VERSION} $*

clean:
	rm -rf _build

purge: clean
	rm -rf spks _tailscale

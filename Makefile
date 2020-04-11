TAILSCALE_VERSION := 0.97-45

.PHONY: tailscale-% clean purge

all: tailscale-amd64 tailscale-386 tailscale-arm64 tailscale-arm

tailscale-%:
	@./build-package.sh ${TAILSCALE_VERSION} $*

clean:
	rm -rf _build _tailscale

purge: clean
	rm -rf spks

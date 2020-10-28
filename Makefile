TAILSCALE_VERSION="1.1.540"
SPK_BUILD="1"

.PHONY: tailscale-% clean purge

all: tailscale-amd64 tailscale-386 tailscale-arm64 tailscale-arm

tailscale-%:
	@./build-package.sh ${TAILSCALE_VERSION} $* ${SPK_BUILD}

clean:
	rm -rf _build _tailscale

purge: clean
	rm -rf spks

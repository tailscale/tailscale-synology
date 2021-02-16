TAILSCALE_VERSION="1.4.4"
TAILSCALE_TRACK="stable"
SPK_BUILD="1"

.PHONY: tailscale-% clean purge

all: tailscale-amd64 tailscale-386 tailscale-arm64 tailscale-arm

tailscale-%:
	@./build-package.sh ${TAILSCALE_TRACK} ${TAILSCALE_VERSION} $* ${SPK_BUILD}

clean:
	rm -rf _build _tailscale

purge: clean
	rm -rf spks

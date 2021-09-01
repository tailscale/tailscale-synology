TAILSCALE_VERSION ?= "1.14.0"
TAILSCALE_TRACK = "stable"
# This needs to be monotinically increasing regardless of the TAILSCALE_VERSION
SPK_BUILD = "006"

.PHONY: tailscale-% clean purge

all: tailscale-amd64 tailscale-386 tailscale-arm64 tailscale-arm

tailscale-%:
	@./build-package.sh ${TAILSCALE_TRACK} ${TAILSCALE_VERSION} $* ${SPK_BUILD} "6"
	@./build-package.sh ${TAILSCALE_TRACK} ${TAILSCALE_VERSION} $* ${SPK_BUILD} "7"

clean:
	rm -rf _build _tailscale

purge: clean
	rm -rf spks

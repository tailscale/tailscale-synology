TAILSCALE_VERSION ?= "1.16.0"
TAILSCALE_TRACK = "stable"
# This needs to be monotinically increasing regardless of the TAILSCALE_VERSION
SPK_BUILD = 11

.PHONY: tailscale-% clean purge

all: tailscale-amd64 tailscale-386 tailscale-arm64 tailscale-arm

release: tailscale-release-amd64 tailscale-release-386 tailscale-release-arm64 tailscale-release-arm

RELEASE_VERSION_ARG="true"
SIDELOAD_VERSION_ARG="false"

tailscale-release-%:
	@./build-package.sh ${TAILSCALE_TRACK} ${TAILSCALE_VERSION} $* ${SPK_BUILD} "6" ${RELEASE_VERSION_ARG}
	@./build-package.sh ${TAILSCALE_TRACK} ${TAILSCALE_VERSION} $* ${SPK_BUILD} "7" ${RELEASE_VERSION_ARG}

tailscale-%:
	@./build-package.sh ${TAILSCALE_TRACK} ${TAILSCALE_VERSION} $* ${SPK_BUILD} "6" ${SIDELOAD_VERSION_ARG}
	@./build-package.sh ${TAILSCALE_TRACK} ${TAILSCALE_VERSION} $* ${SPK_BUILD} "7" ${SIDELOAD_VERSION_ARG}

clean:
	rm -rf _build _tailscale

purge: clean
	rm -rf spks

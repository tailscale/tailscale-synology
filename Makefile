TAILSCALE_VERSION ?= "1.20.1"
TAILSCALE_TRACK ?= "stable"
# SPK_BUILD is derived from the TAILSCALE_VERSION using the forumla
#   (major-1)*1e7 + minor*1e4 + patch*1e1 + dsm
# e.g. for 1.20.1 (dsm7) => 200017
# Note: The DSM version is appended by the build step.
SPK_BUILD ?= 20001

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

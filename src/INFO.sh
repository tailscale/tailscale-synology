#!/bin/bash

VERSION=$1
ARCH=$2
PKG_SIZE=$3
DSM_VERSION=$4

TIMESTAMP=$(date -u +%Y%m%d-%H:%M:%S)

# architecture taken from:
# https://github.com/SynoCommunity/spksrc/wiki/Synology-and-SynoCommunity-Package-Architectures
# https://github.com/SynologyOpenSource/pkgscripts-ng/tree/master/include platform.<PLATFORM> files
case $ARCH in
amd64)
  PLATFORMS="x64 x86 apollolake avoton braswell broadwell broadwellnk bromolow cedarview denverton dockerx64 geminilake grantley purley kvmx64 v1000 x86_64"
  ;;
386)
  PLATFORMS="evansport"
  ;;
arm64)
  PLATFORMS="aarch64 armv8 rtd1296 armada37xx"
  ;;
arm)
  PLATFORMS_ARM5="armv5 88f6281 88f628x"
  PLATFORMS_ARM7="armv7 alpine armada370 armada375 armada38x armadaxp comcerto2k monaco hi3535 ipq806x northstarplus dakota"
  PLATFORMS="${PLATFORMS_ARM5} ${PLATFORMS_ARM7}"
  ;;
*)
  # PLATFORMS_PPC="powerpc ppc824x ppc853x ppc854x qoriq"
  echo "Unsupported architecture: ${ARCH}"
  exit 1
  ;;
esac

if [ "$DSM_VERSION" = "6" ]; then
  os_min_ver="6.0.1-7445"
  os_max_ver="7.0-40000"
else
  os_min_ver="7.0-40000"
  os_max_ver=""
fi

cat <<EOF
package="Tailscale"
version="${VERSION}"
arch="${PLATFORMS}"
description="Connect all your devices using WireGuard, without the hassle."
displayname="Tailscale"
maintainer="Tailscale, Inc."
maintainer_url="https://github.com/tailscale/tailscale-synology"
create_time="${TIMESTAMP}"
dsmuidir="ui"
dsmappname="SYNO.SDS.Tailscale"
startstop_restart_services="nginx"
os_min_ver="${os_min_ver}"
os_max_ver="${os_max_ver}"
extractsize="${PKG_SIZE}"
EOF

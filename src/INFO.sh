#!/bin/bash

VERSION=$1
ARCH=$2
PKG_SIZE=$3
PLATFORMS="apollolake"
TIMESTAMP=$(date -u +%Y%m%d-%H:%M:%S)

cat <<EOF
package="tailscale"
version="${VERSION}"
arch="${PLATFORMS}"
description="Connect all your devices using WireGuard, without the hassle."
displayname="Tailscale"
maintainer="nirev"
maintainer_url="https://github.com/nirev/synology-tailscale"
create_time="${TIMESTAMP}"
extractsize=${PKG_SIZE}
EOF

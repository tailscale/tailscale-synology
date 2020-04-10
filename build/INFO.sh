#!/bin/bash
source /pkgscripts-ng/include/pkg_util.sh

package="tailscale"
version="0.97-45"
displayname="Tailscale"
maintainer="nirev"
arch="$(pkg_get_platform)"
description="Connect all your devices using WireGuard, without the hassle."
[ "$(caller)" != "0 NULL" ] && return 0
pkg_dump_info

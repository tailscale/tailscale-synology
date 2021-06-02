#!/bin/bash

set -e

TAILSCALE_TRACK=$1
TAILSCALE_VERSION=$2
ARCH=$3
SPK_BUILD=$4
DSM_VERSION=$5

download_tailscale() {
  local base_url="https://pkgs.tailscale.com/${TAILSCALE_TRACK}"
  local pkg_name="tailscale_${TAILSCALE_VERSION}_${ARCH}.tgz"
  local src_pkg="${base_url}/${pkg_name}"
  local dest_pkg="_tailscale/${pkg_name}"
  mkdir -p _tailscale

  echo ">>> Downloading package: ${src_pkg}"
  wget --no-verbose -c ${src_pkg} -O ${dest_pkg}

  echo ">>> Extracting.."
  tar -xzf ${dest_pkg} -C _tailscale
}

make_inner_pkg() {
  local tmp_dir=$1
  local dest_dir=$2
  local dest_pkg="$dest_dir/package.tgz"
  local tailscale_dir="_tailscale/tailscale_${TAILSCALE_VERSION}_${ARCH}"

  echo ">>> Making inner package.tgz"
  mkdir -p "${tmp_dir}/bin"
  cp -a ${tailscale_dir}/tailscale{,d} "${tmp_dir}/bin/"

  mkdir -p "${tmp_dir}/ui"
  cp -a src/config "${tmp_dir}/ui/"
  cp -a src/PACKAGE_ICON_256.PNG "${tmp_dir}/ui/"
  cp "${tailscale_dir}/tailscale" "${tmp_dir}/ui/index.cgi"

  mkdir -p "${tmp_dir}/conf"
  cp -a src/tailscaled_logrotate "${tmp_dir}/conf/logrotate.conf"

  pkg_size=$(du -sk "${tmp_dir}" | awk '{print $1}')
  echo "${pkg_size}" >>"$dest_dir/extractsize_tmp"

  ls --color=no "$tmp_dir" | tar -cJf $dest_pkg -C "$tmp_dir" -T /dev/stdin
}

make_spk() {
  local spk_tmp_dir=$1
  local spk_version="${TAILSCALE_VERSION}-${SPK_BUILD}"
  local spk_dest_dir="./spks"
  local pkg_size=$(cat ${spk_tmp_dir}/extractsize_tmp)
  local spk_filename="tailscale-${ARCH}-${spk_version}-dsm${DSM_VERSION}.spk"

  echo ">>> Making spk: ${spk_filename}"
  mkdir -p ${spk_dest_dir}
  rm "${spk_tmp_dir}/extractsize_tmp"

  # copy scripts and icon
  cp -ra src/scripts $spk_tmp_dir
  cp -a src/PACKAGE_ICON*.PNG $spk_tmp_dir
  mkdir ${spk_tmp_dir}/conf
  cp -a "src/privilege-dsm${DSM_VERSION}" ${spk_tmp_dir}/conf/privilege

  cp -a src/Tailscale.sc ${spk_tmp_dir}/Tailscale.sc

  # Generate INFO file
  ./src/INFO.sh "${spk_version}" ${ARCH} ${pkg_size} "${DSM_VERSION}" >"${spk_tmp_dir}"/INFO

  tar -cf "${spk_dest_dir}/${spk_filename}" -C "${spk_tmp_dir}" $(ls ${spk_tmp_dir})
}

make_pkg() {
  mkdir -p ./_build
  local pkg_temp_dir=$(mktemp -d -p ./_build)
  local spk_temp_dir=$(mktemp -d -p ./_build)

  make_inner_pkg ${pkg_temp_dir} ${spk_temp_dir}
  make_spk ${spk_temp_dir}
  echo ">> Done"
  echo ""
  rm -rf ${spk_temp_dir} ${pkg_temp_dir}
}

main() {
  echo "> Building package for TAILSCALE_VERSION=${TAILSCALE_VERSION} ARCH=${ARCH} DSM=${DSM_VERSION}"
  download_tailscale
  make_pkg
}

main

#!/bin/bash

set -e

TAILSCALE_VERSION=$1
ARCH=$2
SPK_BUILD=$3

download_tailscale() {
  local base_url="https://pkgs.tailscale.com/stable"
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
  mkdir -p ${tmp_dir}/bin
  cp -a ${tailscale_dir}/tailscale{,d} ${tmp_dir}/bin/

  mkdir -p ${tmp_dir}/conf
  cp -a src/tailscaled_logrotate ${tmp_dir}/conf/logrotate.conf

  pkg_size=$(du -sk "${tmp_dir}" | awk '{print $1}')
  echo "${pkg_size}" >>"$dest_dir/extractsize_tmp"

  ls --color=no $tmp_dir | tar -cJf $dest_pkg -C "$tmp_dir" -T /dev/stdin
}

spk_build_part() {
  [[ "${SPK_BUILD}" -gt "1" ]] && echo "_spkbuild${SPK_BUILD}"
}

make_spk() {
  local spk_tmp_dir=$1
  local spk_build_part=$(spk_build_part)
  local spk_dest_dir="./spks"
  local pkg_size=$(cat ${spk_tmp_dir}/extractsize_tmp)
  local spk_filename="tailscale_${TAILSCALE_VERSION}${spk_build_part}_${ARCH}.spk"

  echo ">>> Making spk: ${spk_filename}"
  mkdir -p ${spk_dest_dir}
  rm "${spk_tmp_dir}/extractsize_tmp"

  # copy scripts and icon
  cp -ra src/scripts $spk_tmp_dir
  cp -a src/PACKAGE_ICON*.PNG $spk_tmp_dir

  # Generate INFO file
  ./src/INFO.sh "${TAILSCALE_VERSION}${spk_build_part}" ${ARCH} ${pkg_size} >"${spk_tmp_dir}"/INFO

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
  echo "> Building package for TAILSCALE_VERSION=${TAILSCALE_VERSION} ARCH=${ARCH}"
  download_tailscale
  make_pkg
}

main

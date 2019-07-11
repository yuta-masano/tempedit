#!/bin/bash
#===============================================================================
#         USAGE: archive.sh ALL_OS ALL_ARCH PKG_DEST_DIR
#
#   DESCRIPTION: Archive GO binaries generated by each OS and architecture.
#===============================================================================
# Fail on unset variables, command errors and pipe fail.
set -o nounset -o errexit -o pipefail

# Prevent commands misbehaving due to locale differences.
export LC_ALL=C LANG=C

#===============================================================================
#  GLOBAL DECLARATIONS
#===============================================================================
# Script arguments
ALL_OS="$1"
ALL_ARCH="$2"
PKG_DEST_DIR="$3"

#===============================================================================
#  MAIN SCRIPT
#===============================================================================
cd "$PKG_DEST_DIR"
for os in $ALL_OS; do
	for arch in $ALL_ARCH; do
		if $(echo "${os}_${arch}" | grep --quiet 'linux'); then
			tar --create --file="../${os}_${arch}.tar.gz" \
				--auto-compress --verbose "${os}_${arch}"
		else
			zip --recurse-paths "../${os}_${arch}.zip" "${os}_${arch}"
		fi
	done;
done

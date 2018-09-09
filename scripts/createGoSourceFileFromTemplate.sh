#!/bin/bash

set -euo pipefail

BASEDIR="$(cd "$(dirname "$0")/.." && pwd)"

sed \
  -e "s:<cflags>:$(pkg-config --cflags python3):" \
  -e "s:<ldflags>:$(pkg-config --static --libs python3):" \
  "${BASEDIR}/internal/pycgo/updateYAML.go.template" >"${BASEDIR}/internal/pycgo/updateYAML.go"

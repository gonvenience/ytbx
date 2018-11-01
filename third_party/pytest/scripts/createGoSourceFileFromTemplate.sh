#!/bin/bash

set -euo pipefail

BASEDIR="$(cd "$(dirname "$0")/.." && pwd)"

# In case the local Python library is available
if LIB_FILE="$(ls "${BASEDIR}"/third_party/lib/python/lib/libpython3*.a)"; then
  LIB_NAME="$(basename "${LIB_FILE}" | sed -e 's/^lib//' -e 's/.a$//')"
  INCLUDE_PATH="${BASEDIR}/third_party/lib/python/include/${LIB_NAME}"

  sed \
    -e "s:<cflags>:-O3 -pthread -I${INCLUDE_PATH}:" \
    -e "s:<ldflags>:-L${BASEDIR}/third_party/lib/python/lib -l${LIB_NAME} -lm -ldl -lutil:" \
    "${BASEDIR}/internal/pycgo/updateYAML.go.template" >"${BASEDIR}/internal/pycgo/updateYAML.go"

  exit 0
fi

# Default case, use pkg-config
sed \
  -e "s:<cflags>:$(pkg-config --cflags python3):" \
  -e "s:<ldflags>:$(pkg-config --static --libs python3):" \
  "${BASEDIR}/internal/pycgo/updateYAML.go.template" >"${BASEDIR}/internal/pycgo/updateYAML.go"

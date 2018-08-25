#!/bin/bash

set -euo pipefail

BASEDIR="$(cd "$(dirname "$0")/.." && pwd)"

FIND_PYTHON_H="$( (find /usr -type f -name Python.h 2>/dev/null || true) | grep python3)"
if [[ "$(wc -l <<<"${FIND_PYTHON_H}")" -eq 0 ]]; then
  echo "Unable to find Python header file in /usr"
  exit 1
fi

if [[ "$(wc -l <<<"${FIND_PYTHON_H}")" -gt 1 ]]; then
  echo "Found more than one Python header file, not sure what to do ..."
  exit 1
fi

case "$(uname)" in
  Linux)
    FIND_PYTHON_LIB="$( (find /usr -name "libpython3*.so" -exec readlink -f {} \; 2>/dev/null || true) | uniq)"
    if [[ "$(wc -l <<<"${FIND_PYTHON_LIB}")" -eq 0 ]]; then
      echo "Unable to find Python library in /usr"
      exit 1
    fi

    if [[ "$(wc -l <<<"${FIND_PYTHON_LIB}")" -gt 1 ]]; then
      echo "Found more than one Python library, not sure what to do ..."
      exit 1
    fi
    ;;

  \
    Darwin)
    FIND_PYTHON_LIB="$( (find /usr -name "libpython3.?.dylib" 2>/dev/null || true) | grep "lib/lib")"
    if [[ "$(wc -l <<<"${FIND_PYTHON_LIB}")" -eq 0 ]]; then
      echo "Unable to find Python library in /usr"
      exit 1
    fi

    if [[ "$(wc -l <<<"${FIND_PYTHON_LIB}")" -gt 1 ]]; then
      echo "Found more than one Python library, not sure what to do ..."
      exit 1
    fi
    ;;

esac

PATH_TO_PYTHON_H_DIR="$(dirname "${FIND_PYTHON_H}")"
PATH_TO_PYTHON_LIB="$(dirname "${FIND_PYTHON_LIB}")"
PYTHON_LIB_NAME="$(basename "${FIND_PYTHON_LIB}" | sed -e 's/^lib//' -e 's/.so.*//' -e 's/.dylib//')"

echo "Path to Python.h: ${PATH_TO_PYTHON_H_DIR}"
echo "Path to Python library: ${PATH_TO_PYTHON_LIB}"
echo "Name of Python library used: ${PYTHON_LIB_NAME}"

sed \
  -e "s:<path-to-dir-that-has-python-header-file>:${PATH_TO_PYTHON_H_DIR}:" \
  -e "s:<path-to-dir-that-has-python-lib>:${PATH_TO_PYTHON_LIB}:" \
  -e "s:<name-of-python-lib>:${PYTHON_LIB_NAME}:" \
  "${BASEDIR}/internal/pycgo/updateYAML.go.template" >"${BASEDIR}/internal/pycgo/updateYAML.go"

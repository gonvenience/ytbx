#!/bin/bash

set -euo pipefail

BASEDIR="$(cd "$(dirname "$0")/.." && pwd)"

setupDarwin() {
  if [[ ! -d "${BASEDIR}/third_party/lib/python" ]]; then
    mkdir -p "${BASEDIR}/third_party/src"
    pushd "${BASEDIR}/third_party/src" >/dev/null

    curl --silent --location https://www.python.org/ftp/python/3.7.0/Python-3.7.0.tgz | tar -xzf -
    pushd Python-3.7.0 >/dev/null

    ./configure --prefix "${BASEDIR}/third_party/lib/python" --disable-shared --with-openssl="$(brew --prefix openssl)" --enable-optimizations
    make --jobs
    make install
    popd >/dev/null

    popd >/dev/null
  fi
}

setupLinux() {
  if [[ ! -d "${BASEDIR}/third_party/lib/python" ]]; then
    mkdir -p "${BASEDIR}/third_party/src"
    pushd "${BASEDIR}/third_party/src" >/dev/null

    curl --silent --location https://www.python.org/ftp/python/3.6.6/Python-3.6.6.tar.xz | tar -xJf -
    pushd Python-3.6.6 >/dev/null

    ./configure --prefix "${BASEDIR}/third_party/lib/python" --disable-shared --enable-optimizations
    make --jobs
    make install
    popd >/dev/null

    popd >/dev/null
  fi
}

case "$(uname)" in
  Darwin)
    setupDarwin
    ;;

  Linux)
    setupLinux
    ;;
esac

if ! (pip3 list | grep ruamel >/dev/null); then
  "${BASEDIR}"/third_party/lib/python/bin/pip3 install --user --upgrade \
    'pip' \
    'setuptools' \
    'wheel' \
    'ruamel.yaml<=0.15.42' \
    'cython'
fi

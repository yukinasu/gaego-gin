#! /bin/sh -eux

ROOT_DIR="$(cd $(dirname ${BASH_SOURCE:-$0}); pwd)/.."

cd "$ROOT_DIR/server/src/app"

goapp serve

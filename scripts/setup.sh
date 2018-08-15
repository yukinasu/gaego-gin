#! /bin/sh -eux

ROOT_DIR="$(cd $(dirname ${BASH_SOURCE:-$0}); pwd)/.."

cd "$ROOT_DIR/server/src"

# go get
go get -u github.com/golang/dep/cmd/dep
go get -u github.com/golang/lint/golint
go get -u golang.org/x/tools/cmd/goimports
go get -u github.com/kisielk/errcheck
go get -u github.com/swaggo/swag/cmd/swag

# vendor
dep ensure

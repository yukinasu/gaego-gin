#!/bin/sh -eux

ROOT_DIR="$(cd $(dirname ${BASH_SOURCE:-$0}); pwd)/.."

cd "$ROOT_DIR/server/src"

# import文の整理とフォーマッターの適用
goimports -w .

# コードの静的チェック
go vet ./...

# lintの実行
go list ./... | xargs golint -set_exit_status

# エラーハンドリング漏れをチェック
# https://github.com/kisielk/errcheck
errcheck -ignoretests -ignore 'Close' ./...

# swaggoの実行
# https://github.com/swaggo
swag init -g "app/main.go"

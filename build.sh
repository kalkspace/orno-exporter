#!/bin/sh

set -e

LDFLAGS="-X github.com/kalkspace/orno-exporter/config.Version=${GITHUB_SHA} $GO_LDFLAGS"
GO_OPTS="-ldflags \"${LDFLAGS}\""
if [ -z "${BIN_FILE}" ]; then
    GO_OPTS="-o \"${BIN_FILE}\" $GO_OPTS"
fi

go build ${GO_OPTS} main.go

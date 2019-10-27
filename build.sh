#!/bin/sh

set -x

LDFLAGS="-X github.com/kalkspace/orno-exporter/config.Version=${GITHUB_SHA} $GO_LDFLAGS"
GO_OPTS=""
if [ ! -z "${BIN_FILE}" ]; then
    GO_OPTS="-o ${BIN_FILE}"
fi

go build ${GO_OPTS} -ldflags "${LDFLAGS}" main.go

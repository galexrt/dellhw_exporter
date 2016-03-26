#!/bin/bash

dir=$( dirname $0 )

if [ -z ${VERSION} ]; then
    VERSION=0
fi

go get ./...
mkdir ${dir}/dist
GOOS=linux GOARCH=amd64 go build --ldflags "-s -w -X main.BuildDate=`date -u '+%Y-%m-%d_%H:%M'` \
  -X main.HWEVersion=${VERSION}" \
  -o ${dir}/dist/dellhw_exporter

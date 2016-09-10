#!/bin/sh
set -ex

export GOPATH=${PWD}/go

src=${PWD}/go/src/github.com/desource/acbuild-gstore-resource
out=${PWD}/out

CGO_ENABLED=0 go build -o ${out}/gstore ${src}/main.go

cp ${src}/override.sh ${out}/override.sh

cat <<EOF > ${out}/Dockerfile
FROM quay.io/desource/acbuild-resource

COPY gstore       /opt/bin/gstore
COPY override.sh  /opt/resource/override.sh
EOF

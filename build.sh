#!/bin/sh
set -ex

export GOPATH=${PWD}/go

src=${PWD}/go/src/github.com/desource/acbuild-gstore-resource
out=${PWD}/out

go build -o ${out}/gstore ${src}/main.go

cp ${src}/upload.sh ${out}/upload.sh

cat <<EOF > ${out}/Dockerfile
FROM quay.io/desource/acbuild-resource

COPY gstore     /opt/bin/gstore
COPY upload.sh  /opt/resource/upload.sh
EOF

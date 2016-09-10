#!/bin/sh
set -ex

src=${PWD}/src
out=${PWD}/out

go build -o ${out}/gstore ${src}/main.go

cp ${src}/upload.sh ${out}/upload.sh

cat <<EOF > ${out}/Dockerfile
FROM quay.io/desource/acbuild-resource

COPY gstore     /opt/bin/gstore
COPY upload.sh  /opt/resource/upload.sh
EOF

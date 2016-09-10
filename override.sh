#!/bin/sh

upload() {
    destination=$1
    payload=$2

    bucket=$(jq -r '.source.bucket // ""' < $payload)
    bucketPrefix=$(jq -r '.source.bucketPrefix // ""' < $payload)
    email=$(jq -r '.source.email // ""' < $payload)

    private_key_file=$(mktemp /tmp/private-key-file.XXXXXX)

    echo $(jq -r '.source.private_key // ""' < $payload) > $private_key_file

    build=$(jq -r '.params.build // ""' < $payload)
    pwd=$(jq -r '.params.pwd // ""' < $payload)

    if [  -n "$pwd" ]; then
        cd $pwd
    fi

    set -x
    
    /opt/bin/gstore \
        -bucket $bucket \
        -prefix $prefix \
        -email $email \
        -privateKey $private_key_file \
        $build
}

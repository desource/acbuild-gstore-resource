#!/bin/sh

upload() {
    destination=$1
    payload=$2

    bucket=$(jq -r '.source.bucket // ""' < $payload)
    email=$(jq -r '.source.email // ""' < $payload)

    private_key_file=$(mktemp /tmp/private-key-file.XXXXXX)
    
    set -x    
    
    echo $(jq -r '.source.private_key // ""' < $payload) > $private_key_file

    prefix=$(jq -r '.params.prefix // ""' < $payload)
    dir=$(jq -r '.params.dir // ""' < $payload)
    pwd=$(jq -r '.params.pwd // ""' < $payload)

    if [  -n "$pwd" ]; then
        cd $pwd
    fi

    cat $private_key_file

    /opt/bin/gstore \
        -bucket $bucket \
        -prefix $prefix \
        -email $email \
        -privateKey $private_key_file \
        $dir
}

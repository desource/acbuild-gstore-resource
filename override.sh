#!/bin/sh

upload() {
    destination=$1
    payload=$2

    private_key_file=$(mktemp /tmp/private-key-file.XXXXXX)    
    cat <<EOF > $private_key_file
$(jq -r '.source.private_key // ""' < $payload)
EOF
    args="-privateKey=${private_key_file}"
    
    email=$(jq -r '.source.email // ""' < $payload)
    if [ ! -z "${email}" ]; then
        args="${args} -email=${email}"
    fi
    
    bucket=$(jq -r '.source.bucket // ""' < $payload)
    if [ ! -z "${bucket}" ]; then
        args="${args} -bucket=${bucket}"
    fi

    prefix=$(jq -r '.params.prefix // ""' < $payload)
    if [ ! -z "${prefix}" ]; then
        args="${args} -prefix=${prefix}"
    else
        prefix=$(jq -r '.source.prefix // ""' < $payload)
        if [ ! -z "${prefix}" ]; then
            args="${args} -prefix=${prefix}"
        fi
    fi
 
    public=$(jq -r '.params.public // ""' < $payload)
    if [ ! -z "${public}" ]; then
        args="${args} -public=${public}"
    fi
   
    dir=$(jq -r '.params.dir // ""' < $payload)
    if [ ! -z "${dir}" ]; then
        args="${args} ${dir}"
    fi
 
    pwd=$(jq -r '.params.pwd // ""' < $payload)
    if [  -n "$pwd" ]; then
        cd $pwd
    fi

       
    eval /opt/bin/gstore ${args}
}

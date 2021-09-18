#!/usr/bin/env bash

Base="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
case `uname -s` in
    Linux*)     sslConfig=$Base/ssl/openssl.cnf;;
esac
openssl req \
    -newkey rsa:2048 \
    -x509 \
    -nodes \
    -keyout server.key \
    -new \
    -out server.pem \
    -subj /CN=localhost \
    -reqexts SAN \
    -extensions SAN \
    -config <(cat $sslConfig) \
    -sha256 \
    -days 3650

mv $Base/server.key $Base/ssl/server.key
mv $Base/server.pem $Base/ssl/server.pem

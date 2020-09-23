#!/bin/bash

KEYDIR="certs"
cd $KEYDIR

# CA cert
openssl req -nodes -new -x509 -keyout ca.key -out ca.crt -subj "/CN=MiFoMM CA"
# private key
openssl genrsa -out mifomm.key 4096
# Generate and sign the key
openssl req -new -key mifomm.key -subj "/CN=admission-control-service.default.svc." \
    | openssl x509 -req -CA ca.crt -CAkey ca.key -CAcreateserial -out mifomm.crt 
# Create .pem versions
cp mifomm.crt mifommcrt.pem \
    | cp mifomm.key mifommkey.pem
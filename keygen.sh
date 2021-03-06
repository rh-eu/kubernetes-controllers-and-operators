#!/bin/bash

KEYDIR="certs"
mkdir $KEYDIR
cd $KEYDIR

# CA cert
openssl req -nodes -new -x509 -keyout ca.key -out ca.crt -subj "/CN=MiFoMM CA"
# private key
openssl genrsa -out mifomm.key 4096
# Generate and sign the key
openssl req -new -key mifomm.key -subj "/CN=mifomm.validation.svc." \
    | openssl x509 -req -CA ca.crt -CAkey ca.key -CAcreateserial -out mifomm.crt 
# Create .pem versions
cp mifomm.crt mifommcrt.pem \
    | cp mifomm.key mifommkey.pem

# Generate a single line base64 encoded certificate
cat ca.crt | base64 | tr -d '\n' > cabase64.crt
    
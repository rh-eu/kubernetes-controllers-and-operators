# kubernetes-controllers-and-operators

Digging deeeper into admit_funcs.go ...

kind create cluster --name admission-control --config kind/kindv1.19.1.yaml

mkdir -p certs

cd certs
minica -domains mifomm.validation.svc
cd ..

docker build -t mifomm/webhook . && docker push mifomm/webhook

kubectl apply -f YAML/namespace.yaml

kubectl create secret tls mifomm-tls-certs --key="./certs/mifomm.validation.svc/key.pem" --cert="./certs/mifomm.validation.svc/cert.pem" -n validation

kubectl apply -f YAML/deployment/

CA=`cat ./certs/minica.pem | base64 | tr -d '\n'`

cat YAML/webhook.yaml | sed "s/      caBundle: .*/      caBundle: ${CA}/" | kubectl -n validation apply -f -

kubectl -n validation logs -f -l app=validation-webhook

kubectl apply -f samples/enforce-pod-annotations/unannotated-deployment.yaml

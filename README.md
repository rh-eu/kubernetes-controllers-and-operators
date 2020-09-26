# kubernetes-controllers-and-operators

Digging deeper into admit_funcs.go ...

```bash
kind create cluster --name admission-control --config kind/kindv1.19.1.yaml
```

```bash
mkdir -p certs
```

```bash
cd certs
minica -domains mifomm.validation.svc
cd ..
```

```bash
docker build -t mifomm/webhook . && docker push mifomm/webhook
```

```bash
kubectl apply -f YAML/namespace.yaml
```

```bash
kubectl create secret tls mifomm-tls-certs --key="./certs/mifomm.validation.svc/key.pem" --cert="./certs/mifomm.validation.svc/cert.pem" -n validation
```

```bash
kubectl apply -f YAML/deployment/
```

```bash
CA=`cat ./certs/minica.pem | base64 | tr -d '\n'`
```

```bash
cat YAML/webhook.yaml | sed "s/      caBundle: .*/      caBundle: ${CA}/" | kubectl -n validation apply -f -
```

```bash
kubectl -n validation logs -f -l app=validation-webhook
```

```bash
kubectl apply -f samples/enforce-pod-annotations/unannotated-deployment.yaml
```

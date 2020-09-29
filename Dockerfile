FROM golang:1.14.7-alpine3.12 AS build

WORKDIR /go/src/github.com/rh-eu/kubernetes-controllers-and-operators

COPY ./certs/ ./certs/.
COPY go.* ./
COPY ./cmd/ ./cmd/.
COPY ./pkg ./pkg/.


RUN go build -o server cmd/app/main.go

FROM alpine:3.12

USER nobody:nobody
COPY --from=build /go/src/github.com/rh-eu/kubernetes-controllers-and-operators/server /server

CMD [ "/server" ]
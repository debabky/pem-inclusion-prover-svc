FROM golang:1.20-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/debabky/pem-inclusion-prover-svc
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/pem-inclusion-prover-svc /go/src/github.com/debabky/pem-inclusion-prover-svc


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/pem-inclusion-prover-svc /usr/local/bin/pem-inclusion-prover-svc
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["pem-inclusion-prover-svc"]

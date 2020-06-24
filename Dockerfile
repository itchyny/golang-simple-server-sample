FROM golang:1.14-alpine AS builder
RUN apk add --no-cache make git
WORKDIR /workdir
COPY . .
ENV GO111MODULE=on
RUN make build

FROM alpine:3.12
RUN apk add --no-cache ca-certificates
COPY --from=builder /workdir/build/golang-simple-server-sample /
ENTRYPOINT ["/golang-simple-server-sample"]

FROM golang:1.16-alpine AS builder
RUN apk add --no-cache make git
WORKDIR /app
COPY . .
RUN make build

FROM alpine:3.12
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/golang-simple-server-sample /
ENTRYPOINT ["/golang-simple-server-sample"]

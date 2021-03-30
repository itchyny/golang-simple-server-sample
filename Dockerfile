FROM golang:1.16 AS builder

WORKDIR /app
COPY . .
ENV CGO_ENABLED 0
RUN make build

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/golang-simple-server-sample /
ENTRYPOINT ["/golang-simple-server-sample"]

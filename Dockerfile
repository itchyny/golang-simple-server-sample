FROM golang:1.17 AS builder

WORKDIR /app
COPY . .
ENV CGO_ENABLED 0
RUN make build

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/golang-simple-server-sample /
EXPOSE 8080
ENTRYPOINT ["/golang-simple-server-sample"]

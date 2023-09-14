FROM golang:1.21 AS builder

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY *.go Makefile ./
ENV CGO_ENABLED 0
RUN make build

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/golang-simple-server-sample /
EXPOSE 8080
ENTRYPOINT ["/golang-simple-server-sample"]

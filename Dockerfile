FROM --platform=$BUILDPLATFORM golang:1.26 AS builder

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY *.go Makefile ./
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH make build

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/golang-simple-server-sample /
EXPOSE 8080
ENTRYPOINT ["/golang-simple-server-sample"]

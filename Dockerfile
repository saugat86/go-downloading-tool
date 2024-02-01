FROM golang:1.18.1-alpine as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /build


# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o go-downloading-tool

FROM scratch
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/go-downloading-tool .

ENTRYPOINT ["/app/go-downloading-tool"]
CMD ["-h"]

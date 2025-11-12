FROM golang:1.22-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/edge-server ./cmd

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/edge-server /app/edge-server
COPY etc/devices.csv /app/etc/devices.csvs

CMD ["/app/edge-server"]
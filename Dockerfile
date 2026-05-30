FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /deploy-webhook ./cmd/deploy-webhook

FROM scratch
COPY --from=builder /deploy-webhook /deploy-webhook
ENTRYPOINT ["/deploy-webhook"]

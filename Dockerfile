FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git build-base

WORKDIR /app
COPY . .

RUN CGO_ENABLED=1 go build --trimpath -ldflags="-s -w" -o app ./cmd/.

FROM alpine:3.22.1
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8083
CMD ["./app"]

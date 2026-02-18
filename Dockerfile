FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY . .

RUN go build -o app ./cmd/main.go

FROM alpine:3.22.1
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]

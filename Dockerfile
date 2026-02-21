FROM golang:1.24-bookworm AS builder

RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates build-essential \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .

RUN CGO_ENABLED=1 go build --trimpath -ldflags="-s -w" -o app ./cmd/.

FROM debian:bookworm-slim
RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates \
  && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8083
CMD ["./app"]

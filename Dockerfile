FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source and build
COPY . .
RUN go build -ldflags "-X main.version=dev" -o branch-clean

# Minimal runtime image
FROM alpine:3.20

RUN apk add --no-cache git

COPY --from=builder /app/branch-clean /usr/local/bin/branch-clean

ENTRYPOINT ["branch-clean"]

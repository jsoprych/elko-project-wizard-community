# Stage 1: Build
# Pure Go (modernc/sqlite) — no CGO, no gcc needed
FROM golang:1.23-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o elko-project-wizard ./cmd/wizard

# Stage 2: Runtime (minimal)
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /build/elko-project-wizard .
COPY web/ ./web/
COPY directives/ ./directives/

EXPOSE 8080
VOLUME ["/app/data"]

CMD ["./elko-project-wizard", "--port", "8080", "--data", "/app/data"]

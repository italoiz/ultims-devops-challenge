# Build stage
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache ca-certificates

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /api ./cmd/api

# Final stage - scratch-based for smallest possible image (~6MB).
# Using scratch instead of distroless because there are no runtime deps
# and we copy ca-certificates from the build stage.
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /api /api

# Run as nobody (65534) - never run containers as root
USER 65534

EXPOSE 8080 9090

ENTRYPOINT ["/api"]

# Build Stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install Certificates (if needed)
RUN apk --no-cache add ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o logparser

# Final Stage
FROM gcr.io/distroless/static-debian12

# Declare the /data directory as a volume
VOLUME /data

COPY --from=builder /app/logparser /logparser
COPY --from=builder /app/views /views
COPY --from=builder /app/Modules /Modules
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/logparser"]

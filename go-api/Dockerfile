FROM golang:1.24 AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary into /tmp
RUN go build -o /tmp/server ./cmd/server

# ---- Runtime ----
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy only the binary from builder
COPY --from=builder /tmp/server .

EXPOSE 3001

CMD ["./server"]

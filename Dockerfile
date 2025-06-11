# Stage 1: Build
FROM golang:1.24.4-bookworm AS builder

# Install required libraries
RUN apt-get update && apt-get install -y \
    mupdf-tools \
    libmupdf-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .

# Enable CGO (must be 1 for go-fitz to link to libmupdf)
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# Stage 2: Minimal runtime
FROM debian:bookworm-slim

# Install MuPDF CLI tools (needed at runtime)
RUN apt-get update && apt-get install -y \
    mupdf-tools \
    && rm -rf /var/lib/apt/lists/*

# Copy binary from builder
COPY --from=builder /app/main /app/main
COPY --from=builder /app/output /output

WORKDIR /output
CMD ["/app/main"]

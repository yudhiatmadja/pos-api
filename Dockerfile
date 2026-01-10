FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git 
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary
# CGO_ENABLED=0 membuat binary static 
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

# ----------------------------
# STAGE 2: Runner
# ----------------------------
FROM alpine:latest

WORKDIR /root/

# Install CA Certificates 
RUN apk --no-cache add ca-certificates

# Copy binary dari stage builder
COPY --from=builder /app/main .

# Copy folder migrasi 
COPY --from=builder /app/db/migration ./db/migration

# Expose port aplikasi
EXPOSE 8080

CMD ["./main"]
# Stage 1: Build the application
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dvorfs-repository-manager ./cmd/dvorfs-repository-manager

# Stage 2: Create a minimal final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/dvorfs-repository-manager .

EXPOSE 8080

CMD ["./dvorfs-repository-manager"]

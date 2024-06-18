# Build stage
FROM golang:1.22.4-alpine AS builder
WORKDIR /app
COPY ./src/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o clustereye

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/clustereye .
RUN apk add --no-cache curl

CMD ["./clustereye"]
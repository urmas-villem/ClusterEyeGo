# Build Stage
FROM golang:1.22.4-alpine as builder
WORKDIR /work/
COPY ./src/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /clustereye

# Final Stage
FROM golang:1.22.4-alpine
WORKDIR /app
COPY --from=builder /work/clustereye .
RUN apk add --no-cache curl
CMD ["/app/clustereye"]
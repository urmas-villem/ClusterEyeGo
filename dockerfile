
FROM golang:1.22.4-alpine
WORKDIR /app
COPY ./src/ ./
RUN go env -w GO111MODULE=off && \
    CGO_ENABLED=0 GOOS=linux go build -o clustereye
RUN apk add --no-cache curl

CMD ["./clustereye"]
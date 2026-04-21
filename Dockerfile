# Stage 1: build
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /iollama ./cmd/iollama

# Stage 2: final
FROM alpine:latest
RUN apk --no-cache add libstdc++
COPY --from=builder /iollama /usr/local/bin/iollama
EXPOSE 8080
ENTRYPOINT ["iollama"]
CMD ["serve"]

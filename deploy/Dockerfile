# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN  go build cmd/server/main.go

# Runtime stage  
FROM alpine:latest

COPY --from=builder /app .

COPY ../migrations /app/migrations

COPY ../config .

EXPOSE 8080

CMD ["./main"]
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o patisserie ./cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/patisserie .

RUN chmod +x /app/patisserie

EXPOSE 8080

CMD ["./patisserie"]
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o qflow ./main.go

FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/qflow /app/qflow

EXPOSE 3000

CMD ["/app/qflow"]

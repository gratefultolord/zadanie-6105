FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

ENV GOPROXY=https://proxy.golang.org

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .

RUN go build -o main ./cmd/app/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]

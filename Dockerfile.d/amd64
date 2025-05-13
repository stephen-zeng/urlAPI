FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apk update && apk add --no-cache gcc musl-dev
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o urlAPI .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/urlAPI .

EXPOSE 2233

CMD ["./urlAPI"]

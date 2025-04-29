FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.nju.edu.cn/g' /etc/apk/repositories
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o urlAPI .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/urlAPI .

EXPOSE 2233

CMD ["./urlAPI"]

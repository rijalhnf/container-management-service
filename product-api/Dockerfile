FROM golang:1.26.4-alpine AS builder

WORKDIR /app    

COPY go.mod ./
COPY go.sum* ./
RUN go mod download

COPY . .
RUN go build -o product-api

EXPOSE 8888

CMD ["./product-api"]
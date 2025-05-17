FROM golang:1.23.8-alpine AS builder

COPY . /github.com/n1cedek/microservices_course/week2/grpc/source/
WORKDIR /github.com/n1cedek/microservices_course/week2/grpc/source/

RUN go mod download
RUN go build -o ./bin/crud_server cmd/grpc.server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/n1cedek/microservices_course/week2/grpc/source/bin/crud_server .

CMD ["./crud_server"]
FROM golang:1.19.7-alpine3.17 AS builder
WORKDIR /build
ENV GOPROXY https://goproxy.cn,direct
COPY . .
RUN go mod tidy
RUN go mod download
RUN GOARCH=amd64 GOOS=linux go build -o ndor ./cmd

FROM golang:1.19.7-alpine3.17
WORKDIR /
COPY --from=builder /build/ndor .
COPY ./static ./static
CMD ["./ndor"]

FROM golang:1.16.5 AS builder
RUN mkdir /build
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo cmd/cli/main.go

FROM alpine:3.14
RUN mkdir /app
WORKDIR /app
COPY --from=builder /build/main /app/main
COPY config.json /app/config.json
COPY output /app/output

CMD ["/app/main"]
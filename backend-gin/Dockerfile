FROM golang:1.25 AS builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server ./cmd/server

FROM golang:1.25
WORKDIR /app
COPY --from=builder /server /server
EXPOSE 8080
CMD ["/server"]

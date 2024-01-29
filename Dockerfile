FROM golang:1.17-alpine as builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o oats ./cmd

FROM busybox as runner
COPY --from=builder /app/oats /app
ENTRYPOINT ["/app"]
FROM golang:1.13 as build

RUN go env -w GOPROXY=https://goproxy.cn,direct

WORKDIR /go/cache
ADD . .
RUN go mod download


RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o news_watch_notice cmd/news_watch_notice.go
FROM alpine:3.9 as prod
COPY --from=build /go/cache/news_watch_notice /
CMD ["/news_watch_notice"]
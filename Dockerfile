FROM golang:1.12-alpine

LABEL maintainer="hanyajun0123@gmail.com"

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/news_watch_notice

COPY . .
# Download all the dependencies
# https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations
#RUN go get -d -v ./...
RUN pwd

RUN ls
# Install the package and create test binary
RUN go build -o dist/news_watch_notice cmd/news_watch_notice.go && \
    CGO_ENABLED=0 GOOS=linux go test -c
COPY dist/news_watch_notice /usr/bin/
# Run the executable
CMD ["news_watch_notice"]
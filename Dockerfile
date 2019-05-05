FROM golang:1.12-alpine

LABEL maintainer="hanyajun0123@gmail.com"

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/Han-Ya-Jun/news_watch_notice

# Download all the dependencies
# https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations
#RUN go get -d -v ./...
RUN pwd

RUN ls
# Install the package and create test binary
RUN go install cmd/news_watch_notice.go

# Run the executable
CMD ["news_watch_notice"]
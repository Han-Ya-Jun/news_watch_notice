FROM golang:1.12-alpine

LABEL maintainer="hanyajun0123@gmail.com"

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/Han-Ya-Jun/news_watch_notice
COPY . .

# Download all the dependencies
# https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations
#RUN go get -d -v ./...

# Install the package and create test binary
RUN go install -v $GOPATH/src/github.com/Han-Ya-Jun/news_watch_notice/cmd/news_watch_notice.go -o $GOPATH/src/github.com/Han-Ya-Jun/dist/news_watch_notice

# Run the executable
CMD ["news_watch_notice"]
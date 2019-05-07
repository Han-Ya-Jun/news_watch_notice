FROM alpine:3.6

MAINTAINER hanyajun0123@gmail.com
RUN apk update && apk add curl bash tree tzdata \
    && cp -r -f /usr/share/zoneinfo/Hongkong /etc/localtime
ADD news_watch_notice /usr/bin/
ADD news_watch_notice.sha /usr/bin/
CMD ["news_watch_notice"]
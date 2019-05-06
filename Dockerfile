FROM alpine:3.6

MAINTAINER hanyajun0123@gmail.com
RUN  echo 'http://mirrors.ustc.edu.cn/alpine/v3.5/main' > /etc/apk/repositories \
    && echo 'http://mirrors.ustc.edu.cn/alpine/v3.5/community' >>/etc/apk/repositories \
&& apk update && apk add tzdata \
&& ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/Shanghai" > /etc/timezone
ADD news_watch_notice /usr/bin/
ADD news_watch_notice.sha /usr/bin/
CMD ["news_watch_notice"]
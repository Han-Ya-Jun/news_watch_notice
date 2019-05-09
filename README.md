#### 爬取[GoCn每日新闻并推送到微信](https://gocn.vip/explore/category-14)
![image](http://cdn.hanyajun.com/news_watch.png)
#### 使用方法
##### 通过微信通知
```shell
docker run -v /etc/localtime:/etc/localtime:ro \
-e NOTICE_WECHAT_USERS=特鲁尼克 hanyajun/news_watch_notice
```

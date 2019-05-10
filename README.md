#### 爬取[GoCn每日新闻](https://gocn.vip/explore/category-14)并推送到微信/邮箱
![image](http://cdn.hanyajun.com/news_watch.png)
#### 使用方法
##### 通过微信通知

```
docker run -v /etc/localtime:/etc/localtime:ro \
-e NOTICE_WECHAT_USERS=特鲁尼克 hanyajun/news_watch_notice
```
- NOTICE_WECHAT_USERS 代表你要通知的好友的昵称，其中自己微信的文件助手是默认会发的，如果有多个好友，以逗号","隔开。
![image](http://cdn.hanyajun.com/news_notice_wechat.png)
###### 效果
![image](http://cdn.hanyajun.com/wechat.png)

**微信通知有个缺点，就是网页版微信只能有一个终端登录**
##### 通过邮箱推送

```
docker run -v /etc/localtime:/etc/localtime:ro -e NOTICE_TYPE=mail \ //采用邮箱通知，不填则默认微信
-e NOTICE_MAIL_TO=1581532052@qq.com,hanyajun5876@163.com \ //发送
-e NOTICE_MAIL_PWD=******* \ //邮箱smtp授权密码
-e NOTICE_MAIL_PORT=25 \  //smtp端口
-e NOTICE_MAIL_HOST=smtp.qq.com \ //smtp服务器地址
-e NOTICE_MAIL_EMAIL=1581532052@qq.com \ //发送邮箱
-e NOTICE_MAIL_CC=1581532052@qq.com,hanyajun5876@163.com //发送抄送 hanyajun/news_watch_notice
```
###### 效果
![image](http://cdn.hanyajun.com/wechat4.png)



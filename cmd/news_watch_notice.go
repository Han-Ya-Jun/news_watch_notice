package main

import (
	"fmt"
	"github.com/news_watch_notice/pkg/reptile"
	"github.com/news_watch_notice/pkg/wechat"
	"time"
)

func main() {
	/* 登陆微信 */
	err, loginMap := wechat.WechatLogin()
	if err != nil {
		fmt.Printf("login wechat err:%v", err)
		return
	}
	t := time.Tick(time.Minute * 30)
	var flag bool
	var dateTime string
	for {
		/* 爬虫获取新闻 */
		var content string
		nowDateTime := time.Now().Format("2006-01-02")
		if !flag && nowDateTime != dateTime {
			err, contentList := reptile.GetNewsContent()
			if err != nil {
				fmt.Printf("get newsList err:%v", err)
				content = fmt.Sprintf("get newsList err:%v", err)
			} else {
				flag = true
				dateTime = time.Now().Format("2006-01-02")
				for _, c := range contentList {
					content = content + c
				}
			}
			/* 推送消息 */
			err = wechat.WechatSendMsg(content, "filehelper", loginMap)
			if err != nil {
				fmt.Printf("send msg err:%v", err)
			}

		}

		<-t
	}

}

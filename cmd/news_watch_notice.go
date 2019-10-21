package main

import (
	"fmt"
	m "itchat4go/model"
	"news_watch_notice/pkg/github"
	"news_watch_notice/pkg/mail"
	"news_watch_notice/pkg/reptile"
	"news_watch_notice/pkg/slack"
	"news_watch_notice/pkg/wechat"
	"news_watch_notice/utils"
	"strings"
	"time"
)

func main() {
	noticeType := utils.GetValueFromEnv("NOTICE_TYPE")
	var client *mail.Client
	var loginMap m.LoginMap
	var err error
	var userList []string
	var typeFlag bool
	var slackFlag bool
	var webHookUrl string
	var sendObject mail.SendObject
	var githubPushFlag bool
	var githubToken string
	if noticeType == utils.TYPENOCICEMAIL {
		typeFlag = true
		host := utils.GetValueFromEnv("NOTICE_MAIL_HOST")
		port := utils.GetValueFromEnv("NOTICE_MAIL_PORT")
		email := utils.GetValueFromEnv("NOTICE_MAIL_EMAIL")
		password := utils.GetValueFromEnv("NOTICE_MAIL_PWD")
		client = mail.NewMailClient(host, utils.StrToInt(port), email, password)
		toMails := utils.GetValueFromEnv("NOTICE_MAIL_TO")
		ccMails := utils.GetValueFromEnv("NOTICE_MAIL_CC")
		toMailList := strings.Split(toMails, ",")
		ccMailList := strings.Split(ccMails, ",")
		sendObject = mail.SendObject{
			ToMails:     toMailList,
			CcMails:     ccMailList,
			ContentType: "text/html",
		}
	} else if noticeType == utils.TYPENOCTISLACK {
		typeFlag = true
		slackFlag = true
		webHookUrl = utils.GetValueFromEnv("NOTICE_SLACK_WEB_HOOK_URL")
	} else {
		/* 登陆微信 */
		err, loginMap = wechat.WechatLogin()
		noticeWechatUsers := utils.GetValueFromEnv("NOTICE_WECHAT_USERS")
		u := strings.Split(noticeWechatUsers, ",")
		if err != nil {
			fmt.Printf("login wechat err:%v", err)
			return
		}
		userList = wechat.GetSendUsers(loginMap, u)
	}
	if utils.GetValueFromEnv("GITHUB_PUSH") == utils.GITHUBPUSHFLAG {
		githubPushFlag = true
		githubToken = utils.GetValueFromEnv("GITHUB_TOKEN")
	}
	t := time.Tick(time.Minute * 30)
	var flag bool
	var dateTime string
	for {
		/* 爬虫获取新闻 */
		var content string
		nowDateTime := time.Now().Format("2006-01-02")
		if !flag || nowDateTime != dateTime {
			err, contentList := reptile.GetNewsContent(time.Now())
			if err != nil {
				fmt.Printf("get newsList err:%v", err)
			} else {
				flag = true
				dateTime = time.Now().Format("2006-01-02")
				for _, c := range contentList {
					if typeFlag && !slackFlag {
						c = c + "</br>"
					}
					content = content + c
					fmt.Println(c)
				}
			}
			/* 推送消息 */
			if content != "" {
				if githubPushFlag {
					githubContent := ""
					for _, c := range contentList {
						githubContent = githubContent + "- " + c
					}
					er := github.PushGithub(githubToken, time.Now(), githubContent)
					if er != nil {
						fmt.Printf("push to github err:%v", er.Error())
					}
				}
				if !typeFlag {
					err = wechat.WechatSendMsgs(content, userList, loginMap)
				} else if slackFlag {
					err = slack.SenMsgToSlack(webHookUrl, content)
				} else {
					sendObject.Object = "GOCN每日新闻--" + time.Now().Format("2006-01-02")
					sendObject.Content = content
					fmt.Println(content)
					err = client.SendMail(&sendObject)
				}
				if err != nil {
					fmt.Printf("send msg err:%v", err)
				}
			}

		}
		<-t
	}

}

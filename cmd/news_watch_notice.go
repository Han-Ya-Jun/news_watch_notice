package main

import (
	"fmt"
	m "itchat4go/model"
	"news_watch_notice/pkg/github"
	"news_watch_notice/pkg/mail"
	"news_watch_notice/pkg/reptile"
	"news_watch_notice/pkg/slack"
	"news_watch_notice/pkg/wechat"

	md "github.com/russross/blackfriday"

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

	var gocnDateTime string
	var studyDateTime string
	var totalDateTime string
	var gocnFlag bool
	var studyGolangFlag bool
	var flag bool
	for {
		/* 爬虫获取新闻 */
		var content string
		var studyContent string
		nowDateTime := time.Now().Format("2006-01-02")
		if !flag || totalDateTime != nowDateTime {
			var contentList []string
			if !gocnFlag || gocnDateTime != nowDateTime {
				err, contentList = reptile.GetNewsContent(time.Now())
				if err != nil || len(contentList) == 0 {
					fmt.Printf("get newsList err:%v", err)
					gocnFlag = false
				} else {
					gocnFlag = true
					gocnDateTime = time.Now().Format("2006-01-02")
					for _, c := range contentList {
						if typeFlag && !slackFlag {
							c = c + "</br>"
						}
						content = content + c
						fmt.Println(c)
					}
				}
			}
			if !studyGolangFlag || studyDateTime != nowDateTime {
				err, studyContent = reptile.GetStudyGolangContent(time.Now())
				if err != nil || studyContent == "" {
					studyGolangFlag = false
					fmt.Printf("get newsList err:%v", err)
				} else {
					studyGolangFlag = true
					studyDateTime = time.Now().Format("2006-01-02")
				}
			}
			flag = gocnFlag && studyGolangFlag
			if flag {
				totalDateTime = time.Now().Format("2006-01-02")
			}
			/* 推送消息 */
			if content != "" || studyContent != "" {
				if githubPushFlag {
					if content != "" {
						githubContent := ""
						for _, c := range contentList {
							githubContent = githubContent + "- " + c
						}
						er := github.PushGithub(githubToken, time.Now(), githubContent, "gocn")
						if er != nil {
							fmt.Printf("push to github err:%v", er.Error())
						} else {
							fmt.Printf("push gocn_news_set success\n")
						}
						er = github.PushGithub(githubToken, time.Now(), githubContent, "golang_notes")
						if er != nil {
							fmt.Printf("push to github err:%v", er.Error())
						} else {
							fmt.Printf("push gocn_news to golang_notes success\n")
						}
					}
					if studyContent != "" {
						er := github.PushGithub(githubToken, time.Now(), studyContent, "study_golang")
						if er != nil {
							fmt.Printf("push to github err:%v", er.Error())
						} else {
							fmt.Printf("push to golang_notes success")
						}
						er = github.PushGithub(githubToken, time.Now(), studyContent, "gocn_golang")
						if er != nil {
							fmt.Printf("push to github err:%v", er.Error())
						} else {
							fmt.Printf("push to gocn_golang  success")
						}

					}

				}
				if !typeFlag {
					err = wechat.WechatSendMsgs(content, userList, loginMap)
				} else if slackFlag {
					if content != "" {
						err = slack.SenMsgToSlack(webHookUrl, content, "gocn")
						if err != nil {
							println("push slack gocn  err:%v", err)
						} else {
							println("push  slack  success")
						}
					}
					if studyContent != "" {
						err = slack.SenMsgToSlack(webHookUrl, studyContent, "")
						if err != nil {
							println("push slack  golang  err:%v", err)
						} else {
							println("push  slack golang  success")
						}
					}

				} else {
					if content != "" {
						sendObject.Object = "GOCN每日新闻--" + time.Now().Format("2006-01-02")
						sendObject.Content = content
						fmt.Println(content)
						err = client.SendMail(&sendObject)
						if err != nil {
							fmt.Printf("send mail err:%v", err.Error())
						} else {
							println("send mail success")
						}
					}
					if studyContent != "" {
						sendObject.Object = "go语言中文网-每日资讯--" + time.Now().Format("2006-01-02")
						fmt.Println(studyContent)
						result := md.Run([]byte(studyContent))
						sendObject.Content = string(result)
						err = client.SendMail(&sendObject)
						if err != nil {
							fmt.Printf("send mail err:%v", err.Error())
						} else {
							fmt.Print("send mail success")
						}
					}

				}
				if err != nil {
					fmt.Printf("send msg err:%v", err)
				}
			}

		}
		fmt.Printf("flag:%v,gocnDateTime:%v,studyDateTime:%v,totalDateTime:%v,gocnFlag:%v,studyGolangFlag:%v\n", flag, gocnDateTime, studyDateTime, totalDateTime, gocnFlag, studyGolangFlag)
		<-t
	}

}

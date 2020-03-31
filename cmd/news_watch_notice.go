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

	// notice
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
	} else if noticeType == utils.TYPENOCDISABLE {
		// nothing
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

	// push
	if utils.GetValueFromEnv("GITHUB_PUSH") == utils.GITHUBPUSHFLAG {
		githubPushFlag = true
		githubToken = utils.GetValueFromEnv("GITHUB_TOKEN")
	}
	t := time.Tick(time.Minute * 30)

	var gocnDateTime string
	var studyDateTime string
	var gopherDateTime string
	var totalDateTime string
	var gocnFlag bool
	var studyGolangFlag bool
	var gopherDailyFlag bool
	var flag bool
	for {
		/* 爬虫获取新闻 */
		var content string
		var studyContent string
		var gopherDailyContent string
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
			if !gopherDailyFlag || gopherDateTime != nowDateTime {
				var contentList []string
				err, contentList = reptile.GetGopherDailyContent(time.Now())
				if err != nil || len(contentList) == 0 {
					fmt.Printf("get gopher daily content err:%v", err)
					gopherDailyFlag = false
				} else {
					gopherDailyFlag = true
					gopherDateTime = time.Now().Format("2006-01-02")
					for _, c := range contentList {
						//if typeFlag && !slackFlag {
						//	c = c + "</br>"
						//}
						gopherDailyContent = gopherDailyContent + c
						fmt.Println(c)
					}
				}
			}

			flag = gocnFlag && studyGolangFlag && gopherDailyFlag
			if flag {
				totalDateTime = time.Now().Format("2006-01-02")
			}
			/* 推送消息 */
			if content != "" || studyContent != "" || gopherDailyContent != "" {
				if githubPushFlag {
					if content != "" {
						githubContent := ""
						for _, c := range contentList {
							githubContent = githubContent + "- " + c
						}
						er := github.PushGithub(githubToken, time.Now(), githubContent, "gocn_news", "gocn")
						if er != nil {
							fmt.Printf("push to github err:%v", er.Error())
						} else {
							fmt.Printf("push gocn_news_set success\n")
						}
						er = github.PushGithub(githubToken, time.Now(), githubContent, "gocn_news", "golang_notes")
						if er != nil {
							fmt.Printf("push to github err:%v", er.Error())
						} else {
							fmt.Printf("push gocn_news to golang_notes success\n")
						}
					}
					if studyContent != "" {
						er := github.PushGithub(githubToken, time.Now(), studyContent, "go语言中文网(每日资讯)", "gocn")
						if er != nil {
							fmt.Printf("push to github err:%v", er.Error())
						} else {
							fmt.Printf("push to golang_notes success")
						}
						er = github.PushGithub(githubToken, time.Now(), studyContent, "go语言中文网(每日资讯)", "golang_notes")
						if er != nil {
							fmt.Printf("push to github err:%v", er.Error())
						} else {
							fmt.Printf("push to gocn_golang  success")
						}

					}

					if gopherDailyContent != "" {
						er := github.PushGithub(githubToken, time.Now(), gopherDailyContent, "gopherDaily", "gocn")
						if er != nil {
							fmt.Printf("push to github err:%v", er.Error())
						} else {
							fmt.Printf("push to golang_notes success")
						}
						er = github.PushGithub(githubToken, time.Now(), gopherDailyContent, "gopherDaily", "golang_notes")
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
						err = slack.SenMsgToSlack(webHookUrl, studyContent, "goStudy")
						if err != nil {
							println("push slack  golang  err:%v", err)
						} else {
							println("push  slack golang  success")
						}
					}
					if gopherDailyContent != "" {
						err = slack.SenMsgToSlack(webHookUrl, gopherDailyContent, "gopherDaily")
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
					if gopherDailyContent != "" {
						sendObject.Object = "gopherDaily--" + time.Now().Format("2006-01-02")
						fmt.Println(gopherDailyContent)
						result := md.Run([]byte(gopherDailyContent))
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
		fmt.Printf("flag:%v,gocnDateTime:%v,studyDateTime:%v,gopherDateTime:%v,totalDateTime:%v,gocnFlag:%v,studyGolangFlag:%v,gopherDailyFlag:%v\n", flag, gocnDateTime, studyDateTime, gopherDateTime, totalDateTime, gocnFlag, studyGolangFlag, gopherDailyFlag)
		<-t
	}

}

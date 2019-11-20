package slack

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/karriereat/blackfriday-slack"
	"github.com/nlopes/slack"
	bf "github.com/russross/blackfriday/v2"
)

/*
* @Author:hanyajun
* @Date:2019/6/4 0:13
* @Name:slack
* @Function: slack 通知
 */

func SenMsgToSlack(webHookUrl string, content string, from string) error {
	var text string
	var title string
	var authorLink string
	var authorName string
	if from == "gocn" {
		text = content
		title = "<!channel>" + "GOCN每日新闻--" + time.Now().Format("2006-01-02") + " :smile:\n"
		authorLink = "https://gocn.vip/explore/category-14"
		authorName = "GoCn"
	} else {
		renderer := &slackdown.Renderer{}
		md := bf.New(bf.WithRenderer(renderer), bf.WithExtensions(bf.CommonExtensions))
		ast := md.Parse([]byte(content))
		text = string(renderer.Render(ast))
		authorLink = "https://studygolang.com/go/godaily"
		title = "<!channel>" + "go每日资讯--" + time.Now().Format("2006-01-02") + " :smile:\n"
		authorName = "Go语言中文网"
	}
	fmt.Println(text)
	attachment := slack.Attachment{
		Color:         "good",
		Fallback:      "You successfully posted by Incoming Webhook URL!",
		AuthorName:    authorName,
		AuthorSubname: "News",
		AuthorLink:    authorLink,
		AuthorIcon:    "https://gocn.vip/uploads/nav_menu/12.jpg",
		Text:          text,
		Footer:        authorName,
		MarkdownIn:    []string{"text"},
		FooterIcon:    "https://gocn.vip/static/common/avatar-max-img.png",
		Title:         title,
		Ts:            json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}

	msg := slack.WebhookMessage{

		Attachments: []slack.Attachment{attachment},
	}

	err := slack.PostWebhook(webHookUrl, &msg)
	return err
}

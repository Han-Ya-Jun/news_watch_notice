package slack

import (
	"encoding/json"
	"github.com/nlopes/slack"
	"strconv"
	"time"
)

/*
* @Author:hanyajun
* @Date:2019/6/4 0:13
* @Name:slack
* @Function: slack 通知
 */

func SenMsgToSlack(webHookUrl string, content string) error {
	attachment := slack.Attachment{
		Color:         "good",
		Fallback:      "You successfully posted by Incoming Webhook URL!",
		AuthorName:    "GoCn",
		AuthorSubname: "News",
		AuthorLink:    "https://gocn.vip/explore/category-14",
		AuthorIcon:    "https://gocn.vip/uploads/nav_menu/12.jpg",
		Text:          "<!channel>" + "GOCN每日新闻--" + time.Now().Format("2006-01-02") + " :smile:\n" + content,
		Footer:        "Go-Cn",
		FooterIcon:    "https://gocn.vip/static/common/avatar-max-img.png",
		Ts:            json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}
	err := slack.PostWebhook(webHookUrl, &msg)
	return err
}

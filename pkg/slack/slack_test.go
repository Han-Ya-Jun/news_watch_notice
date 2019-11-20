package slack

import (
	"fmt"
	"testing"
	"time"

	"news_watch_notice/pkg/reptile"
)

/*
* @Author:hanyajun
* @Date:2019/11/19 17:50
* @Name:slack
* @Function:
 */
func TestSenMsgToSlack(t *testing.T) {
	_, content := reptile.GetStudyGolangContent(time.Now().Add(time.Hour * -24))
	fmt.Println(content)
	SenMsgToSlack("https://hooks.slack.com/services/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", content, "")
}

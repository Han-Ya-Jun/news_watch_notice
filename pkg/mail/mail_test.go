package mail

import (
	"fmt"
	"news_watch_notice/pkg/reptile"
	"news_watch_notice/utils"
	"testing"
	"time"

	md "github.com/russross/blackfriday"
)

/*
* @Author:hanyajun
* @Date:2019/11/20 16:32
* @Name:mail
* @Function:
 */
func TestClient_SendMail(t *testing.T) {
	client := NewMailClient("smtp.qq.com", utils.StrToInt("465"), "1581532052@qq.com", "vwhfadvqlqwlgcdh")
	sendObject := SendObject{
		ToMails:     []string{"1581532052@qq.com"},
		CcMails:     []string{"1581532052@qq.com"},
		ContentType: "text/html",
	}
	sendObject.Object = "GOCN每日新闻--" + time.Now().Format("2006-01-02")
	_, content := reptile.GetStudyGolangContent(time.Now().Add(time.Hour * -24))
	fmt.Println(content)
	output := md.Run([]byte(content))
	sendObject.Content = string(output)
	err := client.SendMail(&sendObject)
	if err != nil {
		println("send mail err:%v", err.Error())
	} else {
		println("send mail success")
	}
}

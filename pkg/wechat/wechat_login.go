package wechat

import (
	"fmt"
	"github.com/Han-Ya-Jun/qrcode2console"
	"github.com/news_watch_notice/utils"
	e "itchat4go/enum"
	m "itchat4go/model"
	s "itchat4go/service"
	"time"
)

/*
* @Author:15815
* @Date:2019/4/30 0:15
* @Name:login
* @Function:
 */
var (
	uuid string
	//err        error
	//loginMap   m.LoginMap
	//contactMap map[string]m.User
	//groupMap   map[string][]m.User /* 关键字为key的，群组数组 */
)

func WechatLogin() (err error, loginMap m.LoginMap) {

	/* 从微信服务器获取UUID */
	uuid, err = s.GetUUIDFromWX()
	if err != nil {
		return err, m.LoginMap{}
	}

	/* 根据UUID获取二维码 */
	err = s.DownloadImagIntoDir(e.QRCODE_URL+uuid, "./qrcode")
	if err != nil {
		return err, m.LoginMap{}
	}
	//file, _ := os.Open("./qrcode/qrcode.jpg")
	//image,_:=jpeg.Decode(file)
	//输出到控制台
	qr := qrcode2console.NewQRCode2ConsoleWithPath("./qrcode/qrcode.jpg")
	//out,err:=utils.ConvertFromPath("./qrcode/qrcode.jpg",30)
	//if err!=nil{
	//	utils.PanicErr(err)
	//	return err,m.LoginMap{}
	//}
	//utils.Print2Console(out)
	//qr.SetImage(image)
	//panicErr(err)
	//cmd := exec.Command(`cmd`, `/c start ./qrcode/qrcode.jpg`)
	//err = cmd.Run()
	//if err != nil {
	//	return err,m.LoginMap{}
	//}
	qr.Output()

	/* 轮询服务器判断二维码是否扫过暨是否登陆了 */
	for {
		fmt.Println("正在验证登陆... ...")
		status, msg := s.CheckLogin(uuid)
		if status == 200 {
			fmt.Println("登陆成功,处理登陆信息...")
			loginMap, err = s.ProcessLoginInfo(msg)
			if err != nil {
				utils.PanicErr(err)
			}
			fmt.Println("登陆信息处理完毕,正在初始化微信...")
			err = s.InitWX(&loginMap)
			if err != nil {
				utils.PanicErr(err)
			}

			fmt.Println("初始化完毕,通知微信服务器登陆状态变更...")
			err = s.NotifyStatus(&loginMap)
			if err != nil {
				utils.PanicErr(err)
			}
			fmt.Println("通知完毕,本次登陆信息：")
			fmt.Println(e.SKey + "\t\t" + loginMap.BaseRequest.SKey)
			fmt.Println(e.PassTicket + "\t\t" + loginMap.PassTicket)
			break
		} else if status == 201 {
			fmt.Println("请在手机上确认")
		} else if status == 408 {
			fmt.Println("请扫描二维码")
		} else {
			fmt.Println(msg)
		}
	}

	return nil, loginMap

}

func WechatSendMsg(content, toUsername string, loginMap m.LoginMap) error {
	wxSendMsg := m.WxSendMsg{}
	wxSendMsg.Type = 1
	wxSendMsg.Content = content
	wxSendMsg.FromUserName = loginMap.SelfUserName
	wxSendMsg.ToUserName = toUsername
	wxSendMsg.LocalID = fmt.Sprintf("%d", time.Now().Unix())
	wxSendMsg.ClientMsgId = wxSendMsg.LocalID
	fmt.Printf("**********%v", wxSendMsg)
	err := s.SendMsg(&loginMap, wxSendMsg)
	if err != nil {
		fmt.Printf("send msg err:%v", err)
	}
	return err
}

package mail

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

/*
* @Author:15815
* @Date:2019/5/8 17:47
* @Name:mail
* @Function:邮件发送
 */

type Client struct {
	Host     string
	Port     int
	Mail     string
	Password string
}

type SendObject struct {
	ToMails     []string
	CcMails     []string
	Object      string
	ContentType string
	Content     string
}

func NewMailClient(host string, port int, sendMail, password string) *Client {
	return &Client{
		Host:     host,
		Port:     port,
		Mail:     sendMail,
		Password: password,
	}
}
func (m *Client) SendMail(s *SendObject) error {
	mgs := gomail.NewMessage()
	mgs.SetHeader("From", m.Mail)
	mgs.SetHeader("To", s.ToMails...)
	mgs.SetHeader("Cc", s.CcMails...)
	mgs.SetHeader("Subject", s.Object)
	mgs.SetBody(s.ContentType, s.Content)
	d := gomail.NewDialer(m.Host, m.Port, m.Mail, m.Password)
	if err := d.DialAndSend(mgs); err != nil {
		fmt.Printf("send mail err:%v", err)
		return err
	}
	return nil
}

package github

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

/*
* @Author:hanyajun
* @Date:2019/6/21 14:27
* @Name:github
* @Function: 抓取推送到github上去
 */
func PushGithub(token string, publish time.Time, contentList string, from string) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	var rep string
	var path string
	var sepTitle string
	var sep string
	var title string
	var c string
	if from == "gocn" {
		rep = "gocn_news_set"
		path = "README.md"
		sepTitle = "#"
		title = " gocn_news_"
		sep = "##"
		c = "add gocn news--" + publish.Format("2006-01-02")
	} else if from == "golang_notes" {
		c = "add gocn news--" + publish.Format("2006-01-02")
		rep = "golang-notes"
		path = "gocn_news_" + fmt.Sprintf("%d", time.Now().Year()) + ".md"
		sepTitle = "#"
		sep = "##"
		title = " gocn_news_"
	} else if from == "gocn_golang" {
		rep = "gocn_news_set"
		path = "README.md"
		sepTitle = "#"
		sep = "##"
		c = "add go语言中文网(每日资讯)_" + publish.Format("2006-01-02")
		title = " go语言中文网(每日资讯)_"
	} else {
		rep = "golang-notes"
		path = "gocn_news_" + fmt.Sprintf("%d", time.Now().Year()) + ".md"
		sepTitle = "#"
		sep = "##"
		c = "add go语言中文网(每日资讯)_" + publish.Format("2006-01-02")
		title = " go语言中文网(每日资讯)_"
	}

	sha := ""
	content := &github.RepositoryContentFileOptions{
		Message: &c,
		SHA:     &sha,
		Committer: &github.CommitAuthor{
			Name:  github.String("hanyajun"),
			Email: github.String("1581532052@qq.com"),
			Login: github.String("Han-Ya-Jun"),
		},
		Author: &github.CommitAuthor{
			Name:  github.String("hanyajun"),
			Email: github.String("1581532052@qq.com"),
			Login: github.String("Han-Ya-Jun"),
		},
		Branch: github.String("master"),
	}

	op := &github.RepositoryContentGetOptions{}
	repo, _, _, er := client.Repositories.GetContents(ctx, "Han-Ya-Jun", rep, path, op)
	if er != nil || repo == nil {
		fmt.Println(er)
		return er
	}
	content.SHA = repo.SHA
	decodeBytes, err := base64.StdEncoding.DecodeString(*repo.Content)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	oldContentList := strings.Split(string(decodeBytes), sepTitle+" gocn_news_set_2020")
	content.Content = []byte(oldContentList[0] + sepTitle + " gocn_news_set_2020" + "\n" + sep + title + publish.Format("2006-01-02") + "\n" + contentList + "\n" + oldContentList[1])
	_, _, err = client.Repositories.UpdateFile(ctx, "Han-Ya-Jun", rep, path, content)
	if err != nil {
		println(err)
		return err
	}
	return nil

}

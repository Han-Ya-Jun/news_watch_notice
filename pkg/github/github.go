package github

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"strings"
	"time"
)

/*
* @Author:hanyajun
* @Date:2019/6/21 14:27
* @Name:github
* @Function: 抓取推送到github上去
 */
func PushGithub(token string, publish time.Time, contentList string) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	c := "add gocn news--" + publish.Format("2006-01-02")
	sha := ""
	content := &github.RepositoryContentFileOptions{
		Message: &c,
		SHA:     &sha,
		Committer: &github.CommitAuthor{
			Date:  &publish,
			Name:  github.String("hanyajun"),
			Email: github.String("1581532052@qq.com"),
			Login: github.String("Han-Ya-Jun"),
		},
		Author: &github.CommitAuthor{
			Date:  &publish,
			Name:  github.String("hanyajun"),
			Email: github.String("1581532052@qq.com"),
			Login: github.String("Han-Ya-Jun"),
		},
		Branch: github.String("master"),
	}
	op := &github.RepositoryContentGetOptions{}
	repo, _, _, er := client.Repositories.GetContents(ctx, "Han-Ya-Jun", "gocn_news_set", "README.md", op)
	if er != nil || repo == nil {
		fmt.Println(er)
	}
	content.SHA = repo.SHA
	decodeBytes, err := base64.StdEncoding.DecodeString(*repo.Content)
	if err != nil {
		log.Fatalln(err)
	}
	oldContentList := strings.Split(string(decodeBytes), "## gocn_news_set_2019")
	content.Content = []byte(oldContentList[0] + "\n" + "## gocn_news_set_2019" + "\n" + "### gocn_news_" + publish.Format("2006-01-02") + "\n" + contentList + "\n" + oldContentList[1])
	repos, _, err := client.Repositories.UpdateFile(ctx, "Han-Ya-Jun", "gocn_news_set", "README.md", content)
	if err != nil {
		println(err)
		return err
	}
	fmt.Println(repos.SHA)
	return nil

}

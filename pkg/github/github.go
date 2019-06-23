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
		return er
	}
	content.SHA = repo.SHA
	decodeBytes, err := base64.StdEncoding.DecodeString(*repo.Content)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	oldContentList := strings.Split(string(decodeBytes), "## gocn_news_set_2019")
	content.Content = []byte(oldContentList[0] + "\n" + "## gocn_news_set_2019" + "\n" + "### gocn_news_" + publish.Format("2006-01-02") + "\n" + contentList + "\n" + oldContentList[1])
	repos, _, err := client.Repositories.UpdateFile(ctx, "Han-Ya-Jun", "gocn_news_set", "README.md", content)
	fmt.Println(repos.SHA)
	if err != nil {
		println(err)
		return err
	}
	//获取ref
	ref, _, err := client.Git.GetRef(ctx, "Han-Ya-Jun", "gocn_news_set", "heads/master")
	if err != nil && ref == nil {
		println(err)
		return err
	}
	//获取commit
	com, _, err := client.Git.GetCommit(ctx, "Han-Ya-Jun", "gocn_news_set", *ref.Object.SHA)
	if err != nil && com == nil {
		println(err)
		return err
	}
	blob := &github.Blob{
		Content:  github.String("add gocn news"),
		Encoding: github.String("utf-8"),
	}
	//生成blob
	result, _, err := client.Git.CreateBlob(ctx, "Han-Ya-Jun", "gocn_news_set", blob)
	if err != nil && result == nil {
		println(err)
		return err
	}
	//生成tree
	te := github.TreeEntry{
		SHA:  result.SHA,
		Path: github.String("README.md"),
		Mode: github.String("100644"),
		Type: github.String("blob"),
	}
	t := []github.TreeEntry{te}
	tResult, _, err := client.Git.CreateTree(ctx, "Han-Ya-Jun", "gocn_news_set", *com.Tree.SHA, t)
	if err != nil && tResult == nil {
		println(err)
		return err
	}
	//生成commit
	comm := &github.Commit{
		Message: github.String(c),
		Parents: []github.Commit{
			github.Commit{
				SHA: ref.Object.SHA,
			},
		},
		Tree: &github.Tree{
			SHA: tResult.SHA,
		},
	}
	cResult, _, err := client.Git.CreateCommit(ctx, "Han-Ya-Jun", "gocn_news_set", comm)
	if err != nil && cResult == nil {
		println(err)
		return err
	}
	println(cResult)
	return nil

}

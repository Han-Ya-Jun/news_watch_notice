package reptile

import (
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/news_watch_notice/utils"
	"strings"
	"time"
)

/*
* @Author:15815
* @Date:2019/4/30 0:10
* @Name:reptile
* @Function:
 */

// 爬虫CoNews网页
func GetNewsContent() (e error, content []string) {
	var baseUrl string
	c := colly.NewCollector()
	//t:=time.Now().Add(-time.Hour*time.Duration(24))
	data := time.Now().Format("2006-01-02")
	// Find and visit all links
	c.OnHTML("div > h4 > a", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, data) {
			baseUrl = e.Attr("href")
			fmt.Printf("Link found: %q -> %s\n", e.Text, baseUrl)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	e = c.Visit("https://gocn.vip/question/category-14")

	if e != nil {
		return e, nil
	}
	if baseUrl == "" {
		return errors.New("news note update"), nil
	}
	b := colly.NewCollector()

	// Find and visit all links
	i := 0
	contentList := make([]string, 10)
	b.OnHTML("div.mod-body > div > ol > li", func(e *colly.HTMLElement) {
		contentList[i] = utils.TrimQuotes(fmt.Sprintf("%d. %s\n\n", i+1, e.Text))
		i++
		fmt.Printf("%d:%q\n", i, utils.TrimQuotes(fmt.Sprintf("%d. %s\n\n", i+1, e.Text)))
	})
	b.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	e = b.Visit(baseUrl)
	if e != nil {
		return e, nil
	}
	return nil, contentList

}

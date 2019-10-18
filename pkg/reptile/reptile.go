package reptile

import (
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
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
func GetNewsContent(publishTime time.Time) (e error, content []string) {
	var baseUrl string
	c := colly.NewCollector()
	//t:=time.Now().Add(-time.Hour*time.Duration(24))
	data := publishTime.Format("2006-01-02")
	dateOther := publishTime.Format("2006-01-2")
	// Find and visit all links
	c.OnHTML("div > h4 > a", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, data) {
			baseUrl = e.Attr("href")
			fmt.Printf("Link found: %q -> %s\n", e.Text, baseUrl)
		} else if strings.Contains(e.Text, dateOther) {
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
		return errors.New("news not update"), nil
	}
	b := colly.NewCollector()

	// Find and visit all links
	var contentList []string
	b.OnHTML("div.mod-body > div", func(e *colly.HTMLElement) {
		if e.Text != "" {
			reg := "[a-zA-z]+://[^\\s]*"
			title := "[1-5]\\."
			rm, _ := regexp.Compile(reg)
			title2, _ := regexp.Compile(title)
			fmt.Println(e.Text)
			fmt.Println("***********************************************************")
			matched := title2.FindAllStringSubmatchIndex(e.Text, -1)
			fmt.Println(matched)
			indexList := rm.FindAllStringSubmatchIndex(e.Text, -1)
			fmt.Println(indexList)
			fmt.Println(e.Text[indexList[0][0]:e.Text[indexList[0][1]]])
			index:=strings.Index(e.Text,"编辑:")
			for i, v := range matched {
				if v[0] <=index && i < len(matched)-1 {
					content := e.Text[v[0]:matched[i+1][0]]
					if strings.Contains(content, "编辑:") {
						index := strings.Index(content, "编辑:")
						content = content[:index]
						contentList=append(contentList, content+"\n")
						break
					}
					contentList = append(contentList, content+"\n")
				}
			}
		}
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

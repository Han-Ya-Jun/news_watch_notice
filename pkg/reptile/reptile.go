package reptile

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
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
	c.OnHTML("div.title.media-heading > a", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, data) {
			baseUrl = "https://gocn.vip" + e.Attr("href")
			fmt.Printf("Link found: %q -> %s\n", e.Text, baseUrl)
		} else if strings.Contains(e.Text, dateOther) {
			baseUrl = e.Attr("href")
			fmt.Printf("Link found: %q -> %s\n", e.Text, baseUrl)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	e = c.Visit("https://gocn.vip/topics/node18")

	if e != nil {
		return e, nil
	}
	if baseUrl == "" {
		return errors.New("news not update"), nil
	}
	b := colly.NewCollector()

	// Find and visit all links
	var contentList []string
	b.OnHTML("div.card-body.markdown.markdown-toc > ol ", func(e *colly.HTMLElement) {
		e.ForEach("li", func(i int, e *colly.HTMLElement) {
			contentList = append(contentList, fmt.Sprintf("%v.", i+1)+e.Text+"\n")
		})
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

// 爬虫go语言中文网网页
func GetStudyGolangContent(publishTime time.Time) (e error, content string) {
	var baseUrl string
	c := colly.NewCollector()
	//t:=time.Now().Add(-time.Hour*time.Duration(24))
	data := publishTime.Format("2006-01-02")
	// Find and visit all links
	c.OnHTML("dd > div.title", func(e *colly.HTMLElement) {
		link := e.ChildText("a")
		println(link)
		if strings.Contains(link, data) {
			baseUrl = "https://studygolang.com" + e.ChildAttr("a", "href")
			fmt.Printf("Link found: %q -> %s\n", link, baseUrl)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	e = c.Visit("https://studygolang.com/go/godaily")

	if e != nil {
		return e, ""
	}
	if baseUrl == "" {
		return errors.New("news not update"), ""
	}
	b := colly.NewCollector()

	// Find and visit all links
	b.OnHTML("#wrapper > div > div.row > div.col-md-9.col-sm-6 > div.page > div:nth-child(1) > div:nth-child(2) > div", func(e *colly.HTMLElement) {
		if e.Text != "" {
			content = e.Text
		}
	})
	b.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	e = b.Visit(baseUrl)
	if e != nil {
		return e, ""
	}
	return nil, content
}

//
func GetGopherDailyContent(publishTime time.Time) (e error, content []string) {
	baseUrl := "https://gopher-daily.com/issues/202003/issue-%v.md"
	c := colly.NewCollector()
	//t:=time.Now().Add(-time.Hour*time.Duration(24))
	data := publishTime.Format("20060102")
	url := fmt.Sprintf(baseUrl, data)
	// Find and visit all links
	var contentList []string
	c.OnHTML("body > div.container > div > div.offset-lg-1.col > ol", func(e *colly.HTMLElement) {
		e.ForEach("li", func(i int, e *colly.HTMLElement) {
			fmt.Println(e.Text)
			contentList = append(contentList, fmt.Sprintf("- %v.", i+1)+e.Text+"\n")
		})
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	e = c.Visit(url)
	if e != nil {
		return e, nil
	}
	if len(contentList) == 0 {
		return errors.New("news not update"), nil
	}
	return nil, contentList
}

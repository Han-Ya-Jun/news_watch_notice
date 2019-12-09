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
			contentAll := e.Text
			reg := "[a-zA-z]+://[^\\s]*"
			title := "[1-5]\\."
			var index int
			if strings.Contains(contentAll, "编辑:") {
				index = strings.Index(contentAll, "编辑:")
				contentAll = contentAll[:index]
			}
			rm, _ := regexp.Compile(reg)
			title2, _ := regexp.Compile(title)
			matched := title2.FindAllStringSubmatchIndex(contentAll, -1)
			if len(matched) >= 5 {
				exitMap := make(map[string][]int)
				var matchedNew [][]int
				for _, match := range matched {
					fmt.Println(match)
					fmt.Println(contentAll[match[0]:match[1]])
					if _, ok := exitMap[contentAll[match[0]:match[1]]]; !ok {
						exitMap[contentAll[match[0]:match[1]]] = match
						matchedNew = append(matchedNew, match)
					}
				}
				fmt.Println(matched)
				indexList := rm.FindAllStringSubmatchIndex(contentAll, -1)
				fmt.Println(indexList)
				for i, v := range matchedNew {
					if v[0] <= index {
						if i < len(matchedNew)-1 {
							content := contentAll[v[0]:matchedNew[i+1][0]]
							contentList = append(contentList, content+"\n")
						} else {
							content := contentAll[v[0]:index]
							contentList = append(contentList, content+"\n")
							break
						}
					}
				}

			} else {
				fmt.Println("*********************************")
				e.ForEach("li", func(i int, m *colly.HTMLElement) {
					contentList = append(contentList, fmt.Sprintf("%v.", i+1)+m.Text+"\n")
				})

			}
			for _, content := range contentList {
				fmt.Println(content)
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

func trimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "")

	re, _ = regexp.Compile("\\s")
	src = re.ReplaceAllString(src, "")
	return strings.TrimSpace(src)
}

// 爬虫CoNews网页
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

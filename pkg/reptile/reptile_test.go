package reptile

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

/*
* @Author:hanyajun
* @Date:2019/9/23 21:08
* @Name:reptile
* @Function:
 */

func Test_GetStudyGolangContent(t *testing.T) {
	var result string

	for i := 0; i < 4; i++ {
		if contents, err := ioutil.ReadFile("gocn_news_2019.md"); err == nil {
			//因为contents是[]byte类型，直接转换成string类型后会多一行空格,需要使用strings.Replace替换换行符
			result = strings.Replace(string(contents), "\n", "", 1)
		} else {
			fmt.Println(err)
		}
		publishTime := time.Now().Add(time.Hour * -24 * time.Duration(i))
		timestr := publishTime.Format("2006-01-02")
		_, contents := GetStudyGolangContent(publishTime)
		if strings.Contains(result, timestr) {
			oldContentList := strings.Split(result, "## gocn_news_"+timestr)
			newContent := oldContentList[0] + "\n" + "## go语言中文网(每日资讯)_" + timestr + "\n" + contents + "\n" + "## gocn_news_" + timestr + oldContentList[1]
			err2 := ioutil.WriteFile("gocn_news_2019.md", []byte(newContent), 0666)
			fmt.Println(err2)
		}

	}

}

func Test_GetNewsContent(t *testing.T) {
	publishTime := time.Now()
	_, contents := GetNewsContent(publishTime)
	fmt.Println(contents)

}

func Benchmark_GetNewsContent(b *testing.B) {
	for i := 0; i < 4; i++ {
		b.Log("", i)
	}
}
func Benchmark_GetNewsContent2(b *testing.B) {

}

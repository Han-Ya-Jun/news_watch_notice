package utils

import (
	"os"
	"strconv"
)

/*
* @Author:15815
* @Date:2019/4/30 0:00
* @Name:utils
* @Function: 工具类
 */

const (
	TYPENOCICEMAIL = "mail"
	TYPENOCTISLACK = "slack"
	GITHUBPUSHFLAG = "push"
)

//去掉""
func TrimQuotes(s string) string {
	if len(s) >= 2 {
		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

//获取环境变量某个key的值
func GetValueFromEnv(key string) string {

	return os.Getenv(key)
}

func StrToInt(str string) int {
	int, _ := strconv.Atoi(str)
	return int
}


func GetMin(a,b int)int{
	if a>b{
		return b
	}
	return a
}

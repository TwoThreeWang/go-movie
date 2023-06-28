package easyhttp

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	// https://blog.csdn.net/Richard1997/article/details/110739343
	// 常用的第三方api汇总[获取天气]
	// 1. 随机用户(GET请求, JSON格式)。
	// results表示获取用户个数，也就是数组字段results的长度，默认为1；nat应该是国家简称，好像设置了也没啥用。
	url := "https://api.randomuser.me/?nat=US&results=1"
	h := make(map[string]string)
	h["Content-type"] = "application/json"
	p := make(map[string]interface{})
	bytes, err := Get(url, h, p)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(bytes))
}

func TestPost(t *testing.T) {
	// https://zhuanlan.zhihu.com/p/447151163
	//   {
	//  "body": "千里走单骑",
	//  "title": "关于山东为爱出拳",
	//  "userId": 1,
	//  "id": 101
	// }
	url := "http://jsonplaceholder.typicode.com/posts"
	h := make(map[string]string)
	h["Content-type"] = "application/json"
	p := make(map[string]interface{})
	p["userId"] = 1
	p["title"] = "关于山东为爱出拳"
	p["body"] = "千里走单骑"
	bytes, err := Post(url, h, p)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(bytes))
}

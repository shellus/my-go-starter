package main

import (
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"encoding/json"
)

func main() {

}

// 采集明星列表，返回详情页url列表
func cMingXingList(){
	url := "https://www.houyuantuan.com/mingxing/2/";
	doc, _ := goquery.NewDocument(url)
	doc.Find("body > div.wrapper > div.container > div > div.mod-list > div.hot > ul > li").Map(func(i int, s *goquery.Selection) {
		// :nth-child(1) > a.avatar > img
		// a.name
		name := s.Find("a.name").Text()
		href,_ := s.Find("a.name").Attr("href")
		//avatar_img_src, _ := s.Find("a.avatar > img").Attr("src")
		return json.Marshal(struct {
			name string
			href string
		}{
			name: name,
			href: href,
		})
	})
}

// 采集明星详情页面, 返回相册url列表
func cMingXing(url string) []string{
	doc, _ := goquery.NewDocument(url)

	// body > div.wrapper > div.container > div > div.mod-main > div.modules.pic > ul > li:nth-child(1) > div.cover > a > img
	urls := doc.Find("body > div.wrapper > div.container > div > div.mod-main > div.modules.pic > ul > li").Map(func(i int, s *goquery.Selection) {
		// 相册url
		href,_ := s.Find("div.cover > a").Attr("href")
		return href
	})

	return urls;

}


// 获取相册图片url列表
func cXiangCe(url string) []string{
	doc, _ := goquery.NewDocument(url)

	// body > div.wrapper > div.container > div > div.mod-atlas > div.bd > div > div > ul:nth-child(1) > li:nth-child(1) > div.pic > img
	urls := doc.Find("body > div.wrapper > div.container > div > div.mod-atlas > div.bd > div > div > ul:nth-child(1) > li").Map(func(i int, s *goquery.Selection) {
		// 相册url
		href,_ := s.Find("div.pic > img").Attr("src")
		return href
	})

	return urls;
}
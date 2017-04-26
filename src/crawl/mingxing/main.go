package main

import (
	"github.com/PuerkitoBio/goquery"
	"fmt"
)
type MingXingItem struct {
	name string
	url string
}

func main() {
	c := cMingXingList()
	for  _,mingXingItem := range c {
		fmt.Println(mingXingItem.name)
	}
}

// 采集明星列表，返回详情页url列表
func cMingXingList()(l []MingXingItem){

	url := "https://www.houyuantuan.com/mingxing/2/";
	doc, _ := goquery.NewDocument(url)
	_ = doc.Find("body > div.wrapper > div.container > div > div.mod-list > div.hot > ul > li").Map(func(i int, s *goquery.Selection) string {
		// :nth-child(1) > a.avatar > img
		// a.name
		name := s.Find("a.name").Text()
		href,_ := s.Find("a.name").Attr("href")
		//avatar_img_src, _ := s.Find("a.avatar > img").Attr("src")
		l = append(l, MingXingItem{
			name: name,
			url: href,
		})
		return ""
	})
	return
}

// 采集明星详情页面, 返回相册url列表
//func cMingXing(url string) []string{
//	doc, _ := goquery.NewDocument(url)
//
//	// body > div.wrapper > div.container > div > div.mod-main > div.modules.pic > ul > li:nth-child(1) > div.cover > a > img
//	urls := doc.Find("body > div.wrapper > div.container > div > div.mod-main > div.modules.pic > ul > li").Map(func(i int, s *goquery.Selection) {
//		// 相册url
//		href,_ := s.Find("div.cover > a").Attr("href")
//		return href
//	})
//
//	return urls;
//
//}


// 获取相册图片url列表
//func cXiangCe(url string) []string{
//	doc, _ := goquery.NewDocument(url)
//
//	// body > div.wrapper > div.container > div > div.mod-atlas > div.bd > div > div > ul:nth-child(1) > li:nth-child(1) > div.pic > img
//	urls := doc.Find("body > div.wrapper > div.container > div > div.mod-atlas > div.bd > div > div > ul:nth-child(1) > li").Map(func(i int, s *goquery.Selection) {
//		// 相册url
//		href,_ := s.Find("div.pic > img").Attr("src")
//		return href
//	})
//
//	return urls;
//}
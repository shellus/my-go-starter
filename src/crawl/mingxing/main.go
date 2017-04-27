package main

import (
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"path"
	"net/http"
	"os"
	"io"
	"path/filepath"
)

type MingXingItem struct {
	name string
	url  string
}
type XiangCeItem struct {
	name string
	url  string
}

const storePath = `C:\Users\shellus\Pictures\明星图片`

func main() {
	baseUrl := "https://www.houyuantuan.com"
	c := cMingXingList(baseUrl + "/mingxing/2/")
	for _, mingXingItem := range c {

		fmt.Println(baseUrl + mingXingItem.url)
		xiangces := cMingXing(baseUrl + mingXingItem.url)

		for _, xiangCeItem := range xiangces {

			fmt.Println(baseUrl + xiangCeItem.url)
			u2 := cXiangCe(baseUrl + xiangCeItem.url)

			for _, u3 := range u2 {

				fn := storePath + "\\" + mingXingItem.name + "\\" + xiangCeItem.name + "\\" + path.Base(u3)
				fmt.Println("http:" + u3)
				fmt.Println(fn)

				res, _ := http.Get("http:" + u3)

				os.MkdirAll(filepath.Dir(fn), os.FileMode(777))
				file, _ := os.Create(fn)
				io.Copy(file, res.Body)
			}
		}
	}
}

// 采集明星列表，返回详情页url列表
func cMingXingList(url string) (l []MingXingItem) {

	doc, _ := goquery.NewDocument(url)
	_ = doc.Find("body > div.wrapper > div.container > div > div.mod-list > div.hot > ul > li").Map(func(i int, s *goquery.Selection) string {
		// :nth-child(1) > a.avatar > img
		// a.name
		name := s.Find("a.name").Text()
		href, _ := s.Find("a.name").Attr("href")
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
func cMingXing(url string) (l []XiangCeItem) {
	doc, _ := goquery.NewDocument(url)

	// body > div.wrapper > div.container > div > div.mod-main > div.modules.pic > ul > li:nth-child(1) > div.cover > a > img
	_ = doc.Find("body > div.wrapper > div.container > div > div.mod-main > div.modules.pic > ul > li").Map(func(i int, s *goquery.Selection) string {
		// 相册url
		href, _ := s.Find("div.cover > a").Attr("href")

		name := s.Find("div.cover-title > p > a").Text()
		l = append(l, XiangCeItem{
			name: name,
			url: href,
		})
		return ""
	})

	return

}


// 获取相册图片url列表
func cXiangCe(url string) []string {
	doc, _ := goquery.NewDocument(url)

	// body > div.wrapper > div.container > div > div.mod-atlas > div.bd > div > div > ul:nth-child(1) > li:nth-child(1) > div.pic > img
	urls := doc.Find("body > div.wrapper > div.container > div > div.mod-atlas > div.bd > div > div > ul:nth-child(1) > li").Map(func(i int, s *goquery.Selection) string {
		// 相册url
		href, _ := s.Find("div.pic > img").Attr("src")
		return href
	})

	return urls;
}
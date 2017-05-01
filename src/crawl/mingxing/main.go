package main

import (
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"path"
	"net/http"
	"os"
	"io"
	"path/filepath"
	"os/user"
	"time"
)

// 明星详情页结构
type MingXingItem struct {
	name string
	url  string
}
// 相册结构
type XiangCeItem struct {
	mingXingItem MingXingItem
	name         string
	url          string
}
// 图片结构
type TuPianItem struct {
	xiangCeItem XiangCeItem
	url         string
}

var userPath string
var storePath = `/Pictures/明星图片`
var baseUrl = "https://www.houyuantuan.com"

func main() {

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	userPath = usr.HomeDir
	storePath = userPath + storePath

	go cMingXingList(baseUrl + "/mingxing/2/")

	time.Sleep(time.Hour)
}


// 采集明星列表，返回详情页url列表
func cMingXingList(url string) {

	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}
	_ = doc.Find("body > div.wrapper > div.container > div > div.mod-list > div.hot > ul > li").Map(func(i int, s *goquery.Selection) string {
		name := s.Find("a.name").Text()
		href, _ := s.Find("a.name").Attr("href")

		go cMingXing(MingXingItem{
			name: name,
			url: href,
		})

		return ""
	})

}


// 采集明星详情页面, 返回相册url列表
func cMingXing(mingXingItem MingXingItem) {

	doc, err := goquery.NewDocument(baseUrl + mingXingItem.url)
	if err != nil {
		panic(err)
	}

	_ = doc.Find("body > div.wrapper > div.container > div > div.mod-main > div.modules.pic > ul > li").Map(func(i int, s *goquery.Selection) string {
		href, _ := s.Find("div.cover > a").Attr("href")

		name := s.Find("div.cover-title > p > a").Text()

		go cXiangCe(XiangCeItem{
			mingXingItem: mingXingItem,
			name: name,
			url: href,
		})
		return ""
	})
}


// 获取相册图片url列表
func cXiangCe(xiangCeItem XiangCeItem) {
	doc, err := goquery.NewDocument(baseUrl + xiangCeItem.url)
	if err != nil {
		panic(err)
	}
	doc.Find("body > div.wrapper > div.container > div > div.mod-atlas > div.bd > div > div > ul:nth-child(1) > li").Map(func(i int, s *goquery.Selection) string {
		href, _ := s.Find("div.pic > img").Attr("src")
		go downloader(TuPianItem{
			xiangCeItem: xiangCeItem,
			url: href,
		})
		return ""
	})
}

func downloader(tuPianItem TuPianItem) {
	fn := storePath + "\\" + tuPianItem.xiangCeItem.mingXingItem.name + "\\" + tuPianItem.xiangCeItem.name + "\\" + path.Base(tuPianItem.url)
	fmt.Println("http:" + tuPianItem.url)
	fmt.Println(fn)

	res, err := http.Get("http:" + tuPianItem.url)
	if err != nil {
		panic(err)
	}
	os.MkdirAll(filepath.Dir(fn), os.FileMode(777))
	file, err := os.Create(fn)
	if err != nil {
		panic(err)
	}
	io.Copy(file, res.Body)
}

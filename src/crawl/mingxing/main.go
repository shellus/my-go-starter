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
	mingXingItem *MingXingItem
	name         string
	url          string
}
// 图片结构
type TuPianItem struct {
	xiangCeItem *XiangCeItem
	url         string
}

var userPath string
var storePath = `/Pictures/明星图片`
var baseUrl = "https://www.houyuantuan.com"

var timer_cMingXing = time.NewTicker(1 * time.Second)
var timer_cXiangCe = time.NewTicker(1 * time.Second)
var timer_downloader = time.NewTicker(100 * time.Millisecond)

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
func cMingXingList(url string) (err error) {

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return
	}
	_ = doc.Find("body > div.wrapper > div.container > div > div.mod-list > div.hot > ul > li").Map(func(i int, s *goquery.Selection) string {
		name := s.Find("a.name").Text()
		href, _ := s.Find("a.name").Attr("href")

		<-timer_cMingXing.C
		go cMingXing(MingXingItem{
			name: name,
			url: href,
		})

		return ""
	})
	return

}


// 采集明星详情页面, 返回相册url列表
func cMingXing(mingXingItem MingXingItem) (err error) {

	doc, err := goquery.NewDocument(baseUrl + mingXingItem.url)
	if err != nil {
		return
	}

	_ = doc.Find("body > div.wrapper > div.container > div > div.mod-main > div.modules.pic > ul > li").Map(func(i int, s *goquery.Selection) string {
		href, _ := s.Find("div.cover > a").Attr("href")

		name := s.Find("div.cover-title > p > a").Text()

		<-timer_cXiangCe.C
		go cXiangCe(XiangCeItem{
			mingXingItem: &mingXingItem,
			name: name,
			url: href,
		})
		return ""
	})
	return
}


// 获取相册图片url列表
func cXiangCe(xiangCeItem XiangCeItem) (err error) {
	doc, err := goquery.NewDocument(baseUrl + xiangCeItem.url)
	if err != nil {
		return
	}
	doc.Find("body > div.wrapper > div.container > div > div.mod-atlas > div.bd > div > div > ul:nth-child(1) > li").Map(func(i int, s *goquery.Selection) string {
		href, _ := s.Find("div.pic > img").Attr("src")

		<-timer_downloader.C
		go downloader(TuPianItem{
			xiangCeItem: &xiangCeItem,
			url: href,
		})
		return ""
	})
	return
}

func downloader(tuPianItem TuPianItem) (err error) {
	fn := storePath + "\\" + tuPianItem.xiangCeItem.mingXingItem.name + "\\" + tuPianItem.xiangCeItem.name + "\\" + path.Base(tuPianItem.url)
	fmt.Println("http:" + tuPianItem.url)
	fmt.Println(fn)

	res, err := http.Get("http:" + tuPianItem.url)
	if err != nil {
		return
	}
	os.MkdirAll(filepath.Dir(fn), os.FileMode(644))
	file, err := os.Create(fn)
	if err != nil {
		return
	}
	io.Copy(file, res.Body)

	return
}

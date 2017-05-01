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
	"sync"
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

var rateLimit_cMingXingList = make(chan bool, 5)
var rateLimit_cMingXing = make(chan bool, 5)
var rateLimit_cXiangCe = make(chan bool, 10)
var rateLimit_downloader = make(chan bool, 50)

var exit_sync sync.WaitGroup

func main() {

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	userPath = usr.HomeDir
	storePath = userPath + storePath

	exit_sync.Add(1)
	rateLimit_cMingXingList <- true
	go cMingXingList(baseUrl + "/mingxing/2/")

	exit_sync.Wait()
}


// 采集明星列表，返回详情页url列表
func cMingXingList(url string) (err error) {

	defer func() {
		exit_sync.Add(-1)
		<- rateLimit_cMingXingList
		if err := recover(); err != nil {
			fmt.Print(err)
		}
	}()
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}
	_ = doc.Find("body > div.wrapper > div.container > div > div.mod-list > div.hot > ul > li").Map(func(i int, s *goquery.Selection) string {
		name := s.Find("a.name").Text()
		href, _ := s.Find("a.name").Attr("href")
		exit_sync.Add(1)

		rateLimit_cMingXing <- true
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

	defer func() {
		exit_sync.Add(-1)
		<- rateLimit_cMingXing
		if err := recover(); err != nil {
			fmt.Print(err)
		}
	}()
	doc, err := goquery.NewDocument(baseUrl + mingXingItem.url)
	if err != nil {
		panic(err)
	}

	_ = doc.Find("body > div.wrapper > div.container > div > div.mod-main > div.modules.pic > ul > li").Map(func(i int, s *goquery.Selection) string {
		href, _ := s.Find("div.cover > a").Attr("href")

		name := s.Find("div.cover-title > p > a").Text()
		exit_sync.Add(1)

		rateLimit_cXiangCe <- true
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

	defer func() {
		exit_sync.Add(-1)
		<- rateLimit_cXiangCe
		if err := recover(); err != nil {
			fmt.Print(err)
		}
	}()

	doc, err := goquery.NewDocument(baseUrl + xiangCeItem.url)
	if err != nil {
		panic(err)
	}
	doc.Find("body > div.wrapper > div.container > div > div.mod-atlas > div.bd > div > div > ul:nth-child(1) > li").Map(func(i int, s *goquery.Selection) string {
		href, _ := s.Find("div.pic > img").Attr("src")
		exit_sync.Add(1)

		rateLimit_downloader <- true
		go downloader(TuPianItem{
			xiangCeItem: &xiangCeItem,
			url: href,
		})
		return ""
	})
	return
}

func downloader(tuPianItem TuPianItem) (err error) {

	defer func() {
		exit_sync.Add(-1)
		<- rateLimit_downloader
		if err := recover(); err != nil {
			fmt.Print(err)
		}
	}()

	fn := storePath + "\\" + tuPianItem.xiangCeItem.mingXingItem.name + "\\" + tuPianItem.xiangCeItem.name + "\\" + path.Base(tuPianItem.url)
	fmt.Println("http:" + tuPianItem.url)
	fmt.Println(fn)

	res, err := http.Get("http:" + tuPianItem.url)
	if err != nil {
		panic(err)
	}
	os.MkdirAll(filepath.Dir(fn), os.FileMode(644))
	file, err := os.Create(fn)
	if err != nil {
		panic(err)
	}
	io.Copy(file, res.Body)

	return
}

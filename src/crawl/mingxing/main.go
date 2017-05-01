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

	exit := make(chan bool)

	b := boot()

	o1 := make(chan MingXingItem, 1)

	f1(b, o1)
	f1(b, o1)
	f1(b, o1)
	f1(b, o1)
	f1(b, o1)
	f1(b, o1)
	f1(b, o1)
	f1(b, o1)

	o2 := make(chan XiangCeItem, 1)

	f2(o1, o2)
	f2(o1, o2)
	f2(o1, o2)
	f2(o1, o2)
	f2(o1, o2)
	f2(o1, o2)
	f2(o1, o2)
	f2(o1, o2)
	f2(o1, o2)
	f2(o1, o2)
	f2(o1, o2)
	f2(o1, o2)
	f2(o1, o2)
	f2(o1, o2)
	f2(o1, o2)
	f2(o1, o2)

	o3 := make(chan TuPianItem, 1)

	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)
	f3(o2, o3)

	dc := 15

	for i := 0; i < dc; i++ {
		downloader(o3, exit)
	}

	for i := 0; i < dc; i++ {
		<-exit
	}

}
func boot() chan string {
	o := make(chan string, 1)
	go func() {
		o <- baseUrl + "/mingxing/2/"
		close(o)
	}()
	return o
}

// 列表工厂
// 输入列表url
// 输出详情
func f1(i chan string, o chan MingXingItem) {

	go func() {
		for u := range i {
			for _, item := range cMingXingList(u) {
				o <- item
			}
		}
	}()
}

// 详情工厂
// 输入详情
// 输出相册
func f2(i chan MingXingItem, o chan XiangCeItem) {

	go func() {
		for u := range i {
			for _, item := range cMingXing(u) {
				o <- item
			}
		}
	}()
}

// 相册工厂
// 输入相册
// 输出图片
func f3(i chan XiangCeItem, o chan TuPianItem) {

	go func() {
		for u := range i {
			for _, item := range cXiangCe(u) {
				o <- item
			}
		}
	}()
}

func downloader(i chan TuPianItem, e chan bool) {
	go func() {
		for u := range i {
			fn := storePath + "\\" + u.xiangCeItem.mingXingItem.name + "\\" + u.xiangCeItem.name + "\\" + path.Base(u.url)
			fmt.Println("http:" + u.url)
			fmt.Println(fn)

			res, _ := http.Get("http:" + u.url)

			os.MkdirAll(filepath.Dir(fn), os.FileMode(777))
			file, _ := os.Create(fn)
			io.Copy(file, res.Body)
		}
		e <- true
	}()
}

// 采集明星列表，返回详情页url列表
func cMingXingList(url string) (l []MingXingItem) {

	doc, _ := goquery.NewDocument(url)
	_ = doc.Find("body > div.wrapper > div.container > div > div.mod-list > div.hot > ul > li").Map(func(i int, s *goquery.Selection) string {
		name := s.Find("a.name").Text()
		href, _ := s.Find("a.name").Attr("href")
		l = append(l, MingXingItem{
			name: name,
			url: href,
		})
		return ""
	})
	return
}


// 采集明星详情页面, 返回相册url列表
func cMingXing(mingXingItem MingXingItem) (l []XiangCeItem) {

	doc, _ := goquery.NewDocument(baseUrl + mingXingItem.url)

	_ = doc.Find("body > div.wrapper > div.container > div > div.mod-main > div.modules.pic > ul > li").Map(func(i int, s *goquery.Selection) string {
		href, _ := s.Find("div.cover > a").Attr("href")

		name := s.Find("div.cover-title > p > a").Text()
		l = append(l, XiangCeItem{
			mingXingItem: mingXingItem,
			name: name,
			url: href,
		})
		return ""
	})

	return

}


// 获取相册图片url列表
func cXiangCe(xiangCeItem XiangCeItem) (l []TuPianItem) {
	doc, _ := goquery.NewDocument(baseUrl + xiangCeItem.url)

	doc.Find("body > div.wrapper > div.container > div > div.mod-atlas > div.bd > div > div > ul:nth-child(1) > li").Map(func(i int, s *goquery.Selection) string {
		href, _ := s.Find("div.pic > img").Attr("src")
		l = append(l, TuPianItem{
			xiangCeItem: xiangCeItem,
			url: href,
		})
		return ""
	})

	return
}
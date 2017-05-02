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
	"queue"
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

var q_cMingXingList = queue.NewQueue(5)
var q_cMingXing = queue.NewQueue(5)
var q_cXiangCe = queue.NewQueue(10)
var q_downloader = queue.NewQueue(50)

func main() {

	initPath()

	q_cMingXingList.Sub(func(j queue.Job) {
		cMingXingList(j.Value.(string))
	})

	q_cMingXing.Sub(func(j queue.Job) {
		cMingXing(j.Value.(MingXingItem))
	})
	q_cXiangCe.Sub(func(j queue.Job) {
		cXiangCe(j.Value.(XiangCeItem))
	})
	q_downloader.Sub(func(j queue.Job) {
		downloader(j.Value.(TuPianItem))
	})

	q_cMingXingList.Push(queue.Job{Value:baseUrl + "/mingxing/2/"})

	q_cMingXingList.Work()
	fmt.Println("列表采集完毕")

	q_cMingXing.Work()
	fmt.Println("详情页采集完毕")

	q_cXiangCe.Work()
	fmt.Println("相册采集完毕")

	q_downloader.Work()
	fmt.Println("图片采集完毕")
}

func initPath() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	userPath = usr.HomeDir
	storePath = userPath + storePath
}


// 采集明星列表
func cMingXingList(url string) {

	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}

	doc.Find("body > div.wrapper > div.container > div > div.mod-list > div.hot > ul > li").
		Map(func(i int, s *goquery.Selection) string {
		name := s.Find("a.name").Text()
		href, _ := s.Find("a.name").Attr("href")

		q_cMingXing.Push(queue.Job{
			Value:
			MingXingItem{
				name: name,
				url: href,
			}})

		return ""
	})
}


// 采集明星详情页面
func cMingXing(mingXingItem MingXingItem) {

	doc, err := goquery.NewDocument(baseUrl + mingXingItem.url)
	if err != nil {
		panic(err)
	}

	doc.Find("body > div.wrapper > div.container > div > div.mod-main > div.modules.pic > ul > li").
		Map(func(i int, s *goquery.Selection) string {
		href, _ := s.Find("div.cover > a").Attr("href")

		name := s.Find("div.cover-title > p > a").Text()

		q_cXiangCe.Push(queue.Job{
			Value:
			XiangCeItem{
				mingXingItem: &mingXingItem,
				name: name,
				url: href,
			}})

		return ""
	})
}


// 获取相册图片url列表
func cXiangCe(xiangCeItem XiangCeItem) {

	doc, err := goquery.NewDocument(baseUrl + xiangCeItem.url)
	if err != nil {
		panic(err)
	}

	doc.Find("body > div.wrapper > div.container > div > div.mod-atlas > div.bd > div > div > ul:nth-child(1) > li").
		Map(func(i int, s *goquery.Selection) string {
		href, _ := s.Find("div.pic > img").Attr("src")

		q_downloader.Push(queue.Job{
			Value:
			TuPianItem{
				xiangCeItem: &xiangCeItem,
				url: href,
			}})

		return ""
	})
}

func downloader(tuPianItem TuPianItem) {

	fn := storePath + "\\" + tuPianItem.xiangCeItem.mingXingItem.name + "\\" + tuPianItem.xiangCeItem.name + "\\" + path.Base(tuPianItem.url)

	err := os.MkdirAll(filepath.Dir(fn), os.FileMode(644))
	if err != nil {
		panic(err)
	}

	fmt.Println("http:" + tuPianItem.url)
	fmt.Println(fn)

	res, err := http.Get("http:" + tuPianItem.url)
	if err != nil {
		panic(err)
	}


	file, err := os.Create(fn)
	if err != nil {
		panic(err)
	}

	io.Copy(file, res.Body)
}

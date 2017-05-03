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
var storePath 	= `/Pictures/明星图片`
var baseUrl 	= "https://www.houyuantuan.com"

var q_cMingXingList 	= queue.NewQueue(5)
var q_cMingXing 	= queue.NewQueue(10)
var q_cXiangCe 		= queue.NewQueue(50)
var q_downloader 	= queue.NewQueue(200)

func main() {

	initPath()

	q_cMingXingList.Sub(cMingXingList)
	q_cMingXing.Sub(cMingXing)
	q_cXiangCe.Sub(cXiangCe)
	q_downloader.Sub(downloader)

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

// 获取保存目录
func initPath() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	userPath = usr.HomeDir
	storePath = userPath + storePath
}


// 采集明星列表
func cMingXingList(j queue.Job) {
	url := j.Value.(string)

	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}

	doc.Find("body > div.wrapper > div.container > div > div.mod-list > div.hot > ul > li").
		Each(func(_ int, s *goquery.Selection) {
		name := s.Find("a.name").Text()
		href, _ := s.Find("a.name").Attr("href")

		q_cMingXing.Push(queue.Job{
			Value:
			MingXingItem{
				name: name,
				url: href,
			}})
	})
}


// 采集明星详情页面
func cMingXing(j queue.Job) {
	mingXingItem := j.Value.(MingXingItem)

	doc, err := goquery.NewDocument(baseUrl + mingXingItem.url)
	if err != nil {
		panic(err)
	}

	doc.Find("body > div.wrapper > div.container > div > div.mod-main > div.modules.pic > ul > li").
		Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Find("div.cover > a").Attr("href")

		name := s.Find("div.cover-title > p > a").Text()

		q_cXiangCe.Push(queue.Job{
			Value:
			XiangCeItem{
				mingXingItem: &mingXingItem,
				name: name,
				url: href,
			}})

	})
}


// 获取相册图片url列表
func cXiangCe(j queue.Job) {
	xiangCeItem := j.Value.(XiangCeItem)

	doc, err := goquery.NewDocument(baseUrl + xiangCeItem.url)
	if err != nil {
		panic(err)
	}

	doc.Find("body > div.wrapper > div.container > div > div.mod-atlas > div.bd > div > div > ul:nth-child(1) > li").
		Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Find("div.pic > img").Attr("src")

		q_downloader.Push(queue.Job{
			Value:
			TuPianItem{
				xiangCeItem: &xiangCeItem,
				url: href,
			}})

	})
}

// 下载图片文件
func downloader(j queue.Job) {
	tuPianItem := j.Value.(TuPianItem)

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

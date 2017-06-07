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
	"time"
	"strings"
)





func main() {

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

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	userPath = usr.HomeDir
	storePath = userPath + storePath

	var q_cMingXingList = queue.NewQueue(1)
	var q_cMingXing = queue.NewQueue(1)
	var q_cXiangCe = queue.NewQueue(1)
	var q_downloader = queue.NewQueue(1)

	q_cMingXingList.Sub(func (j queue.Job) {
		url := j.Value.(string)

		doc, err := goquery.NewDocument(url)
		if err != nil {
			panic(err)
		}

		doc.Find("body > div.wrapper > div.container > div > div.mod-list > div.hot > ul > li").
			Each(func(_ int, s *goquery.Selection) {
			name := s.Find("a.name").Text()
			href, _ := s.Find("a.name").Attr("href")

			q_cMingXing.Pub(queue.Job{
				Value:
				MingXingItem{
					name: name,
					url: href,
				}})
			fmt.Println(name)
		})
		doc.Find("body > div.wrapper > div.container > div > div.mod-list > div.list > ul > li").
			Each(func(_ int, s *goquery.Selection) {
			name := s.Find("a").Text()
			href, _ := s.Find("a").Attr("href")
			name = strings.TrimSpace(name)

			q_cMingXing.Pub(queue.Job{
				Value:
				MingXingItem{
					name: name,
					url: href,
				}})
			fmt.Println(name)
		})
	})

	q_cMingXing.Sub(func (j queue.Job) {
		mingXingItem := j.Value.(MingXingItem)

		if mingXingItem.name == "李宇春" {

			return
		}

		doc, err := goquery.NewDocument(baseUrl + mingXingItem.url)
		if err != nil {
			panic(err)
		}

		doc.Find("body > div.wrapper > div.container > div > div.mod-main > div.modules.pic > ul > li").
			Each(func(_ int, s *goquery.Selection) {
			href, _ := s.Find("div.cover > a").Attr("href")

			name := s.Find("div.cover-title > p > a").Text()

			q_cXiangCe.Pub(queue.Job{
				Value:
				XiangCeItem{
					mingXingItem: &mingXingItem,
					name: name,
					url: href,
				}})

		})
	})

	q_cXiangCe.Sub(func (j queue.Job) {
		xiangCeItem := j.Value.(XiangCeItem)

		doc, err := goquery.NewDocument(baseUrl + xiangCeItem.url)
		if err != nil {
			panic(err)
		}

		doc.Find("body > div.wrapper > div.container > div > div.mod-atlas > div.bd > div > div > ul:nth-child(1) > li").
			Each(func(_ int, s *goquery.Selection) {
			href, _ := s.Find("div.pic > img").Attr("src")

			q_downloader.Pub(queue.Job{
				Value:
				TuPianItem{
					xiangCeItem: &xiangCeItem,
					url: href,
				}})

		})
	})

	q_downloader.Sub(func (j queue.Job) {
		tuPianItem := j.Value.(TuPianItem)

		fn := storePath + "\\" + tuPianItem.xiangCeItem.mingXingItem.name + "-" + tuPianItem.xiangCeItem.name + "-" + path.Base(tuPianItem.url)
		fn, err := filepath.Abs(fn)
		if err != nil {
			panic(err)
		}
		err = os.MkdirAll(filepath.Dir(fn), os.FileMode(644))
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
	})

	q_cMingXingList.Pub(queue.Job{Value:baseUrl + "/mingxing/2/"})

	q_cMingXingList.Work()
	fmt.Println("列表采集完毕")

	q_cMingXing.Work()
	fmt.Println("详情页采集完毕")

	q_cXiangCe.Work()
	fmt.Println("相册采集完毕")
	time.Sleep(time.Second)

	q_downloader.Work()
	fmt.Println("图片采集完毕")
}


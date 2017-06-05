package main

import (
	"net/http"
	"path/filepath"
	"github.com/hprose/hprose-golang/util"
	"os"
	"io"
	"github.com/pkg/errors"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"queue"
	"golang.org/x/net/html"
	"fmt"
	"sync"
)

const baseDir string = `C:\Users\shellus\Downloads\sex\`
const baseUrl string = "http://99vv1.com"

var c = http.DefaultClient

var qItemPage = queue.NewQueue(5, "ItemPage")

var qVideoDownload = queue.NewQueue(20, "VideoDownload")

func main() {

	qItemPage.Sub(func(j queue.Job) (err error) {
		defer func (){
			if r := recover(); r != nil{
				// todo: 隐瞒错误
				fmt.Println(r)
			}
		}()
		return praseItemPage(j.Value.(string))
	})
	qVideoDownload.Sub(func(j queue.Job) (err error) {
		defer func (){
			if r := recover(); r != nil{
				// todo: 隐瞒错误
				fmt.Println(r)
			}
		}()
		arr := j.Value.([2]string)
		return downloadVideo(arr[0], arr[1])
	})
	//u := "http://99vv1.com/get_file/3/8d4dc711f6a07357db8998b6f334c918/68000/68763/68763.mp4"
	//u := "http://99vv1.com/get_file/3/437fa99544f68479af41cbd43500b40b/68000/68763/68763_hq.mp4"
	//downloadVideo(u)
	//praseItemPage(baseUrl + "videos/68901/2-720p/")


	go qItemPage.Work()
	go qVideoDownload.Work()

	err := parseIndexPage()
	if err != nil {
		fmt.Println(err)
	}

	w := sync.WaitGroup{}
	w.Add(1)
	w.Wait()
}

func parseIndexPage() (err error) {
	req, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {
		return errors.New(err.Error())
	}

	handleRequest(req)

	r, err := c.Do(req)
	if err != nil {
		return errors.New(err.Error())
	}
	defer r.Body.Close()

	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return errors.New(err.Error())
	}
	sel := doc.Find("a")

	for i := 0; i < sel.Length(); i++ {
		ss := sel.Eq(i)

		href, ok := ss.Attr("href")
		if !ok {
			continue
		}
		// href e.g /videos/68649/644806e950cdef46d754e34b4ba04b05/

		match, err := regexp.MatchString("^/videos/\\d*?/.*/$", href)
		if err != nil {
			return errors.New(err.Error())
		}
		if match {
			qItemPage.Push(queue.Job{Value:baseUrl + href})
		}

	}
	return

}

func praseItemPage(u string) (err error) {
	fmt.Println("开始解析页面：" + u)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return errors.New(err.Error())
	}

	handleRequest(req)

	r, err := c.Do(req)
	if err != nil {
		return errors.New(err.Error())
	}
	defer r.Body.Close()

	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return errors.New(err.Error())
	}

	sel := doc.Find("body > div.wrapper > div.content > div > div > div.main > div.wrap-video > div > script").Eq(1)
	name := doc.Find("body > div.wrapper > div.content > div > div > div.main > div:nth-child(1) > div > h1").Text()

	script, err := sel.Html()
	if err != nil {
		return errors.New(err.Error())
	}

	rScript, err := regexp.Compile(`var flashvars = \{([\s\S]*?)\}`)
	if err != nil {
		return errors.New(err.Error())
	}

	json := html.UnescapeString(rScript.FindString(script)[16:])
	/*
	e.g:
{
  video_id: '68901',
  license_code: '$323962310688084',
  lrc: '8d4cb6dd4dc155b3945fb45e5ea53bb0',
  video_url: '/get_file/3/271eae19e43854870ff86e424b13108f/68000/68901/68901.mp4/',
  postfix: '.mp4',
  video_url_text: 'LQ',
  video_alt_url: '/get_file/3/5e26a658e558eaff0f57cae7ba4476a6/68000/68901/68901_hq.mp4/',
  video_alt_url_text: 'HD',
  preview_url: 'http://www.99vv1.com/contents/videos_screenshots/68000/68901/preview.mp4.jpg',
  skin: '1',
  video_click_url: 'http://www.99kk3.com',
  bt: '5',
  hide_controlbar: '0',
  mlogo: '久久热',
  mlogo_link: 'http://www.99kk3.com',
  disable_selected_slot_restoring: 'true',
  adreplay: 'true',
  disable_preview_resize: 'true',
  embed: '1'
}
	 */

	rJson, err := regexp.Compile(`video_url: '.*?'`)
	if err != nil {
		return errors.New(err.Error())
	}
	videoUrl := rJson.FindString(json)
	videoUrl = videoUrl[12:len(videoUrl) - 2]
	/*
	e.g:
_file/3/271eae19e43854870ff86e424b13108f/68000/68901/68901.mp4/
 	*/
	qVideoDownload.Push(queue.Job{Value:[2]string{name, baseUrl + videoUrl}})
	return
}

func downloadVideo(name string, u string) (err error) {
	fmt.Println("开始下载视频：" + u)


	// 请求
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return errors.New(err.Error())
	}

	handleRequest(req)

	r, err := c.Do(req)

	if err != nil {
		return errors.New(err.Error())
	}
	defer r.Body.Close()


	if r.StatusCode != 200 {
		return errors.New("r.StatusCode is " + util.Itoa(r.StatusCode))
	}


	// 创建文件
	f := name + "." + filepath.Ext(u)

	fh, err := os.Create(baseDir + f)

	if err != nil {
		return errors.New(err.Error())
	}
	defer fh.Close()


	// 下载
	_, err = io.Copy(fh, r.Body)

	// 错误回退
	if err != nil {
		os.Remove(baseDir + f)
		return errors.New(fmt.Sprintf("%s \n%s \n%s", err, baseDir + f, u))
	}
	return
}

func handleRequest(req *http.Request) {
	// 必要
	req.Header.Set("Cookie", "PHPSESSID=b62sj4dfpajfffdkatkr5ffnd2; _gat=1; kt_tcookie=1; _ga=GA1.2.175262314.1496382929; _gid=GA1.2.205912476.1496382929; kt_is_visited=1")

	// 非必要
	req.Header.Set("Referer", "http://www.99vv1.com/videos/68763/d35a82d8802e8b8be6ed19ed98a56c64/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
}


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
	"net/url"
	"io/ioutil"
	"bytes"
	"net/http/cookiejar"
	"reflect"
)

const baseDir string = `C:\Users\shellus\Downloads\sex\`
const baseUrl string = "http://99vv1.com"

type Re99VideoInfo struct {
	Url      string
	Title    string
	VideoUrl string
}

var c = http.DefaultClient

var qItemPage = queue.NewQueue(5, "ItemPage")

var qVideoDownload = queue.NewQueue(20, "VideoDownload")

func main() {
	// 持久化cookie
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	c.Jar = jar

	// 模拟登录产生cookie
	err = login()
	if err != nil {
		panic(err)
	}

	// 准备消费job
	qItemPage.Sub(func(j queue.Job) (err error) {
		return praseItemPage(j.Value.(string))
	})
	qVideoDownload.Sub(func(j queue.Job) (err error) {
		maps := j.Value.(map[string]interface{})

		videoInfo := &Re99VideoInfo{}
		structValue := reflect.ValueOf(videoInfo).Elem()

		for k,v := range maps{
			structFieldValue := structValue.FieldByName(k)
			val := reflect.ValueOf(v)
			structFieldValue.Set(val)
		}

		return downloadVideo(videoInfo)
	})

	// 启动工作线程
	go qItemPage.Work()
	go qVideoDownload.Work()

	// 直接访问首页获得第一批种子
	err = parseIndexPage()
	if err != nil {
		panic(err)
	}

	// 防止程序退出
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
			qItemPage.Pub(&queue.Job{Value:baseUrl + href})
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

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.New(err.Error())
	}

	if bytes.Index(buf, []byte("请登陆后观看")) != -1 {
		return errors.New("no login")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(buf))
	if err != nil {
		return errors.New(err.Error())
	}

	sel := doc.Find("body > div.wrapper > div.content > div > div > div.main > div.wrap-video > div > script").Eq(1)
	name := doc.Find("body > div.wrapper > div.content > div > div > div.main > div:nth-child(1) > div > h1").Text()

	script, err := sel.Html()
	script = html.UnescapeString(script)
	if len(script) < 20 {
		return errors.New(fmt.Sprintf("script is empty: %v", script))
	}

	if err != nil {
		return errors.New(err.Error())
	}

	rScript, err := regexp.Compile(`var flashvars = \{([\s\S]*?)\}`)
	if err != nil {
		return errors.New(err.Error())
	}

	rScriptr := rScript.FindString(script)
	ioutil.WriteFile("log.html", buf, os.FileMode(777))
	json := rScriptr[16:]
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

	qVideoDownload.Pub(&queue.Job{
		Value:Re99VideoInfo{
			Url:u,
			Title:name,
			VideoUrl:baseUrl+videoUrl,
		},
	})
	return
}

func downloadVideo(videoInfo *Re99VideoInfo) (err error) {
	fmt.Println("开始下载视频：" + videoInfo.VideoUrl)

	// 请求
	req, err := http.NewRequest("GET", videoInfo.VideoUrl, nil)
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
	f := videoInfo.Title + filepath.Ext(videoInfo.VideoUrl)
	f = baseDir + f

	fh, err := os.Create(f)

	if err != nil {
		return errors.New(err.Error())
	}


	// 下载
	_, err = io.Copy(fh, r.Body)

	fh.Close()
	// 错误回退
	if err != nil {
		os.Remove(f)
		return errors.New(fmt.Sprintf("%s \n%s \n%s", err, f, videoInfo.VideoUrl))
	}
	return
}

func handleRequest(req *http.Request) {
	// 非必要
	//req.Header.Set("Referer", "http://www.99vv1.com/")
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
}

func login() (err error) {
	// 禁止跳转跟随的http客户端

	data := url.Values{
		"action":[]string{"login"},
		"username":[]string{"shellus"},
		"pass":[]string{"a7245810"},
	}

	res, err := c.PostForm("http://www.99vv1.com/login.php", data)
	if err != nil {
		return
	}

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = errors.New(err.Error())
		return
	}

	if s := bytes.Index(buf, []byte(`<div class="message_error">`)); s != -1 {
		s = s + len(`<div class="message_error">`)
		t := bytes.Index(buf[s:], []byte(`</div>`))
		text := "未知错误"
		if t != -1 {
			text = string(buf[s:s + t]);
		}
		err = errors.New(fmt.Sprintf("login error: %s", text))
		return
	}

	return
}
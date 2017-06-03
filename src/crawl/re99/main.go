package main

import (
	"net/http"
	"fmt"
	"path/filepath"
	"github.com/hprose/hprose-golang/util"
	"os"
	"io"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"queue"
	"golang.org/x/net/html"
)

const baseDir string = `C:\Users\shellus\Downloads\sex\`

var c = http.DefaultClient

var qItemPage = queue.NewQueue(1, "ItemPage")

var qVideoDownload = queue.NewQueue(1, "VideoDownload")

func main() {

	//qItemPage.Sub(func(j queue.Job){praseItemPage(j.Value.(string))})

	//u := "http://99vv1.com/get_file/3/8d4dc711f6a07357db8998b6f334c918/68000/68763/68763.mp4"
	//u := "http://99vv1.com/get_file/3/437fa99544f68479af41cbd43500b40b/68000/68763/68763_hq.mp4"
	//downloadVideo(u)
	praseItemPage("http://99vv1.com/videos/68901/2-720p/")
}

func praseIndexPage() {
	u := "http://www.99vv1.com"
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		panic(err)
	}

	// 必要
	req.Header.Set("Cookie", "PHPSESSID=ged6979fa0b4vj0de00ar76vb4; _gat=1; kt_tcookie=1; _ga=GA1.2.175262314.1496382929; _gid=GA1.2.205912476.1496382929; kt_is_visited=1")

	// 非必要
	req.Header.Set("Referer", "http://www.99vv1.com/videos/68763/d35a82d8802e8b8be6ed19ed98a56c64/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	r, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		panic(err)
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
			panic(err)
		}
		if match {
			qItemPage.Push(queue.Job{Value:href})
		}

	}

}

func praseItemPage(u string) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		panic(err)
	}

	// 必要
	req.Header.Set("Cookie", "PHPSESSID=ged6979fa0b4vj0de00ar76vb4; _gat=1; kt_tcookie=1; _ga=GA1.2.175262314.1496382929; _gid=GA1.2.205912476.1496382929; kt_is_visited=1")

	// 非必要
	req.Header.Set("Referer", "http://www.99vv1.com/videos/68763/d35a82d8802e8b8be6ed19ed98a56c64/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	r, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		panic(err)
	}

	sel := doc.Find("body > div.wrapper > div.content > div > div > div.main > div.wrap-video > div > script").Eq(1)

	script, err := sel.Html()
	if err != nil {
		panic(err)
	}

	reg, err := regexp.Compile(`var flashvars = \{([\s\S]*?)\}`);
	if err != nil {
		panic(err)
	}
	fmt.Println(script)

	fmt.Println("aaa")
	fmt.Println(html.UnescapeString(reg.FindString(script)[16:]))
}

func downloadVideo(u string) (l int64) {

	//c := &http.Client{
	//	CheckRedirect: func(req *http.Request, via []*http.Request) error {
	//		return http.ErrUseLastResponse
	//	},
	//}


	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		panic(err)
	}

	// 必要
	req.Header.Set("Cookie", "PHPSESSID=ged6979fa0b4vj0de00ar76vb4; _gat=1; kt_tcookie=1; _ga=GA1.2.175262314.1496382929; _gid=GA1.2.205912476.1496382929; kt_is_visited=1")

	// 非必要
	req.Header.Set("Referer", "http://www.99vv1.com/videos/68763/d35a82d8802e8b8be6ed19ed98a56c64/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	r, err := c.Do(req)

	if err != nil {
		panic(err)
	}

	defer r.Body.Close()

	if r.StatusCode != 200 {
		panic(errors.New("r.StatusCode is " + util.Itoa(r.StatusCode)))
	}

	f := filepath.Base(u)

	if f == "" {
		f = util.UUIDv4()
	}

	fh, err := os.Create(baseDir + f)

	defer fh.Close()

	if err != nil {
		panic(err)
	}

	l, err = io.Copy(fh, r.Body)

	if err != nil {
		panic(err)
	}

	return l
}


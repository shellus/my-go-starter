package example_test

import (
	"fmt"
	"os"
	"container/list"
	"reflect"
	"strconv"
	"net/http"
	"regexp"
	"io/ioutil"
	"bytes"
	"net/url"
	"github.com/antonholmquist/jason"
	"testing"
)

type MyStruct struct {
	name string
}

func TestHttpProxyAccess(t *testing.T){
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://127.0.0.1:1080")
	}

	transport := &http.Transport{Proxy: proxy}

	client := &http.Client{Transport: transport}

	rest, err := client.Get("http://www.tianyancha.com/company/1534045940")
	if err != nil {
		t.Error(err)
	}
	buf , err := ioutil.ReadAll(rest.Body)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(buf))
}

func TestJsonFractal(t *testing.T) {
	str := `{
																	video_id: '71591', 																	license_code: '$316329710436270', 																	lrc: '8d4cb6dd4dc155b3945fb45e5ea53bb0', 																	video_url: '/get_file/3/a06ffe6085ab694adf9b052714f79891/71000/71591/71591.mp4/', 																	postfix: '.mp4', 																	video_url_text: 'LQ', 																	video_alt_url: '/hd.php', 																	video_alt_url_text: 'HD', 																	video_alt_url_redirect: '1', 																	preview_url: 'http://www.99kk5.com/contents/videos_screenshots/71000/71591/preview.mp4.jpg', 																	skin: '1', 																	video_click_url: 'http://www.99kk3.com', 																	bt: '5', 																	hide_controlbar: '0', 																	mlogo: '久久热', 																	mlogo_link: 'http://www.99kk3.com', 																	disable_selected_slot_restoring: 'true', 																	adv_start_html: '/player/html.php?adv_id=start_html&video_id=71591&cs_id=0', 																	adv_pause_html: '/player/html.php?adv_id=pause_html&video_id=71591&cs_id=0', 																	adreplay: 'true', 																	disable_preview_resize: 'true', 																	embed: '1'															}`
	buf := jsonFractal([]byte(str))

	j, err := jason.NewObjectFromBytes(buf)
	if err != nil {
		t.Error(err)
	}
	previewUrl, err := j.GetString("preview_url")
	if err != nil {
		t.Error(err)
	}
	t.Log(previewUrl)
}

func jsonFractal(s []byte)([]byte){
	r,_ := regexp.Compile(`[\w_-]+?: `)
	s = r.ReplaceAllFunc(s, func(s []byte) []byte { return []byte("\"" + string(s[0:len(s)-2]) + "\": ") })
	s = bytes.Replace(s, []byte{39}, []byte{34}, -1)
	s = bytes.Replace(s, []byte{10}, []byte{}, -1)
	s = bytes.Replace(s, []byte{13}, []byte{}, -1)
	s = bytes.Replace(s, []byte{9},  []byte{}, -1)
	s = bytes.Replace(s, []byte{32},  []byte{}, -1)
	return s
}

// 禁止跳转的http客户端
func TestHttpNoCheckRedirect(t *testing.T) {

	c := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := c.Get("http://google.com")
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 301 {
		t.Errorf("http status %d", resp.StatusCode)
	}
}

/**
字符串，数值转换
 */
func TestConvert(t *testing.T) {

	myi, _ := strconv.Atoi("123")
	fmt.Printf("convert to int: %d \n", myi)

	fmt.Printf("convert to int string: %s \n", strconv.Itoa(123))

	buf := []byte{'a', 'b', 'c'}
	fmt.Printf("convert to string: %s \n", string(buf))

	mystring := "abc"
	fmt.Printf("convert to bytes: %s \n", []byte(mystring))

	b, _ := strconv.ParseBool("true")
	f, _ := strconv.ParseFloat("3.1415", 64)
	i, _ := strconv.ParseInt("-42", 10, 64)
	u, _ := strconv.ParseUint("42", 10, 64)

	t.Log(b, f, i, u)
}
func TestReflect(t *testing.T){
	bs := []byte{'a', 'b', 'c'}

	t.Log(reflect.TypeOf(bs).String())
}
/*
链表及字节数组
反射获取类型
 */
func TestMylist(t *testing.T) {

	l := list.New()


	bs := []byte{'a', 'b', 'c'}

	l.PushFront(bs[0])
	l.PushFront(bs[1])
	l.PushFront(bs[2])

	// todo 什么鬼。这个还没搞懂怎么用。

	t.Log(l.Front().Value.(byte))
	t.Log(l.Front().Value.(byte))
	t.Log(l.Front().Value.(byte))
}

/*
动态数组（切片）（slice）
 */
func TestSlice(t *testing.T) {
	lines := []string{"1", "2"}
	lines = append(lines, "3")
	t.Log(lines);
}

/*
字符串连接
 */
func TestStringCombine(t *testing.T) {
	s := "hello"
	s = "c" + s[1:] // 字符串虽不能更改，但可进行切片操作
	t.Logf("string: %s\n", s)
}

/*
数值计算
 */
func TestNumericCalc(t *testing.T) {
	a, b, c := 1, 2, 3
	n := a + b + c;

	t.Logf("calc value: %s\n", n);
}

/*
数组操作
 */
func TestArray(t *testing.T) {
	a := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8}

	_ = a[3:4]

	fmt.Print(len(a));
}

func TestMyMap(t *testing.T) {
	//var numbers map[string]int
	// 另一种map的声明方式
	numbers := make(map[string]int)
	numbers["one"] = 1  //赋值
	numbers["ten"] = 10 //赋值
	numbers["three"] = 3

	t.Logf("第三个数字是: %d", numbers["three"]) // 读取数据

	for k := range numbers {
		t.Log(k)
	}

}

func TestException(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			_ = true
		}
	}()
	ext()
}

func ext() {
	var user = os.Getenv("USER")
	if user == "" {
		panic("no value for $USER")
	}
}

func TestMyType(t *testing.T) {
	type person struct {
		name string
		age  int
	}

	p := person{"abc", 18}

	t.Log(p)

	p2 := person{"def", 19}

	t.Log(p2)

}

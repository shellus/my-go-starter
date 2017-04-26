package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"net/url"
)

func main() {
	str, err := post();

	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Print(string(str));
}


func post()(body []byte, err error){

	my_url:= "http://shcms.localhost/test"

	data := make(url.Values)

	data.Add("username", "shellus")

	resp, err := http.PostForm(my_url, data)
	if(err != nil){return}

	body, err = ioutil.ReadAll(resp.Body)
	if(err != nil){return}

	return
}
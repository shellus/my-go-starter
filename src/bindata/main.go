package main

import (
	"fmt"
)

func main() {

	data, err := Asset("asset/a.txt")
	if err != nil {
		fmt.Print(err)
	}else {
		fmt.Printf("%s", data)
	}
}

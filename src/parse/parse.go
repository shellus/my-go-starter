package main

import (
	"fmt"
	"flag"
	"os"
	"bufio"
	"container/list"
)

func main() {
	filename := flag.String(
		"filename",
		"",
		"input parse filename")

	flag.Parse()

	f, err := os.Open(*filename)

	if err != nil {
		panic(err)
	}

	buf := bufio.NewReader(f)

	lines := list.New()

	for err == nil {
		str, _, err := buf.ReadLine()

		if err != nil {
			break
		}
		if len(str) > 0 && str[0] == '[' {
			lines.PushBack(string(str))
		} else {

			lines.Back().Value = lines.Back().Value.(string) + string(str)
		}
	}
	fmt.Printf("slice length: %d \n", lines.Len())
}

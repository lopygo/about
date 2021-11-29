package main

import (
	"fmt"

	"github.com/lopygo/about/info"
)

func main() {

	i, err := info.FromInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	i.ShowAll()
}

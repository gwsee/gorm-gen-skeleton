package main

import (
	"encoding/json"
	"fmt"
	"github.com/h2non/filetype"
	"net/url"
	"os"
)

type name struct {
	Host *string
}

func main() {
	str := "http://localhost"
	vals, err := url.Parse(str)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		by, _ := json.Marshal(vals)
		fmt.Println(string(by))
	}
	//var a name
	//fmt.Println(a.Host)
	//fmt.Println(*a.Host)
	buf, _ := os.ReadFile("some.1jpg")
	kind, _ := filetype.Match(buf)
	if kind == filetype.Unknown {
		fmt.Println("Unknown file type")
		return
	}

	fmt.Printf("File type: %s. MIME: %s\n", kind.Extension, kind.MIME.Value)
}

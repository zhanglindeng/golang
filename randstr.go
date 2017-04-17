package main

import (
	"github.com/Unknwon/com"
	"fmt"
	"flag"
)

var length = flag.Int("length", 32, "randstr.exe --length=32 default 32")

func main() {
	flag.Parse()
	fmt.Println(string(com.RandomCreateBytes(*length)))
}

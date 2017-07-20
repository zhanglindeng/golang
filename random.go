package main

import (
	"flag"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"strings"

	"github.com/zhanglindeng/phpgo"
)

var port = flag.String("port", "8090", "--port=8090 default 8090")

var (
	s = `0123456789`
	x = `abcdefghijklmnopqrstuvwxyz`
	d = `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
	t = `!@#$%^&*()_-=+|{}[];.,'"<>?~\/`
)

func main() {

	flag.Parse()

	http.HandleFunc("/", index)

	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {

	_type := r.URL.Query().Get("type")
	_length := r.URL.Query().Get("length")

	// 默认 32 个字符
	length, err := strconv.Atoi(_length)
	if err != nil {
		length = 32
	}

	if length <= 0 {
		length = 32
	}

	// 最长 1024 个字符
	if length > 1024 {
		length = 1024
	}

	// sxdt
	if ok, err := regexp.MatchString("^[sxdt]+$", _type); !ok || err != nil {
		_type = "sxd"
	}

	seed := ""
	if strings.Contains(_type, "s") {
		seed = seed + s
	}
	if strings.Contains(_type, "x") {
		seed = seed + x
	}
	if strings.Contains(_type, "d") {
		seed = seed + d
	}
	if strings.Contains(_type, "t") {
		seed = seed + t
	}

	w.Write([]byte(phpgo.StringRandomStr(length, []byte(seed)...)))
}

package main

import (
	"flag"
	"log"
	"net/http"

	"encoding/json"

	"github.com/yanyiwu/gojieba"
)

var fenci *gojieba.Jieba

var port = flag.String("port", "8000", "usage /path/to/fenci --port=8000 default 8000")

func main() {

	flag.Parse()

	fenci = gojieba.NewJieba()
	defer fenci.Free()

	http.HandleFunc("/", index)

	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {

	str := r.URL.Query().Get("s")
	keywords := fenci.CutForSearch(str, true)

	var data = struct {
		S        string   `json:"s"`
		Keywords []string `json:"keywords"`
	}{
		S:        str,
		Keywords: keywords,
	}

	b, _ := json.Marshal(data)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}

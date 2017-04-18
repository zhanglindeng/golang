package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"log"
	"github.com/PuerkitoBio/goquery"
	"time"
	"strconv"
	"net/url"
	"io/ioutil"
	"net/http"
	"encoding/json"
)

type Address struct {
	gorm.Model
	AreaId   int `json:"area_id,string"`
	AreaName string `json:"area_name"`
	ParentId int `json:"parent_id,string"`
	Sort     int `json:"sort,string"`
}

type Brand struct {
	gorm.Model
	Name    string `json:"name"`
	Image   string `json:"image"`
	Address string `json:"address"`
	Trade   string `json:"trade"`
}

func init() {
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Llongfile)
}

func main2() {

	// http://yske.org/index.php?m=union&c=union&a=getArea
	// areaId:130000
	resp, err := http.PostForm("http://yske.org/index.php?m=union&c=union&a=getArea", url.Values{"areaId": {"130000"}})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var cities []Address
	log.Println(string(body))

	err = json.Unmarshal(body, &cities)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println(cities)
}

func main() {
	db, err := gorm.Open("mysql", "root:root@tcp(localhost:33062)/yske?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = db.DB().Ping()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	doc, err := goquery.NewDocument("http://yske.org/index.php?m=union&c=union&a=index&page=1")
	time.Sleep(2 * time.Second)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	li := doc.Find("#province option")
	log.Println(li.Length())
	li.Each(func(i int, s *goquery.Selection) {
		areaId, _ := s.Attr("value")
		areaName := s.Text()
		log.Println(areaId, areaName)

		aid, _ := strconv.Atoi(areaId)
		address := Address{AreaId: aid, AreaName: areaName, Sort: 100 + i, ParentId: 0}
		b := db.NewRecord(address)
		if b {
			log.Println("OK")
			// 获取省下的市区
			fetchCity(db, areaId)
		}
		db.Create(&address)
	})
}

func fetchCity(db *gorm.DB, areaId string) {
	// http://yske.org/index.php?m=union&c=union&a=getArea
	// areaId:130000
	log.Println(areaId)
	resp, err := http.PostForm("http://yske.org/index.php?m=union&c=union&a=getArea", url.Values{"areaId": {areaId}})
	time.Sleep(2 * time.Second)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var cities []Address
	//log.Println(string(body))

	err = json.Unmarshal(body, &cities)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	//log.Println(cities)
	for i, v := range cities {
		address := Address{AreaId: v.AreaId, AreaName: v.AreaName, Sort: v.Sort + i + i, ParentId: v.ParentId}
		b := db.NewRecord(address)
		if b {
			log.Println("OK")
			// 获取市下的市区县市
			fetchArea(db, strconv.Itoa(v.AreaId))
		}
		db.Create(&address)
	}
}

func fetchArea(db *gorm.DB, areaId string) {
	log.Println(areaId)
	resp, err := http.PostForm("http://yske.org/index.php?m=union&c=union&a=getArea", url.Values{"areaId": {areaId}})
	time.Sleep(2 * time.Second)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var cities []Address
	//log.Println(string(body))

	err = json.Unmarshal(body, &cities)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	//log.Println(cities)
	for i, v := range cities {
		address := Address{AreaId: v.AreaId, AreaName: v.AreaName, Sort: v.Sort + i + i, ParentId: v.ParentId}
		b := db.NewRecord(address)
		if b {
			log.Println("OK")
		}
		db.Create(&address)
	}
}

func main1() {
	db, err := gorm.Open("mysql", "root:root@tcp(localhost:33062)/yske?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = db.DB().Ping()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	db.AutoMigrate(&Address{}, &Brand{})

	// http://yske.org/index.php?m=union&c=union&a=index&page=1
	pageTotal := 123

	for counter := 1; counter <= pageTotal; counter++ {
		doc, err := goquery.NewDocument("http://yske.org/index.php?m=union&c=union&a=index&page=" + strconv.Itoa(counter))
		time.Sleep(2 * time.Second)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		li := doc.Find(".sjalllist>ul li")
		log.Println(li.Length())
		li.Each(func(i int, s *goquery.Selection) {
			image, _ := s.Find("img").Attr("src")
			trade := s.Find("span").Text()
			name := s.Find(".sjinfo a").Text()
			address := s.Find(".sjinfo p").Last().Find("i").Text()
			log.Println(image, trade, name, address)

			brand := Brand{Image: image, Trade: trade, Name: name, Address: address}
			b := db.NewRecord(brand)
			if b {
				log.Println("OK")
			}
			db.Create(&brand)
		})
	}
}

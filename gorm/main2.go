package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

var (
	err    error
	db     *gorm.DB
	driver = "mysql"
	source = `root:@tcp(localhost:3306)/gorm?charset=utf8&parseTime=True&loc=Local`
)

type Person struct {
	ID   int64 `json:"id"`
	Name string `json:"name"`
	Age  uint8 `json:"age"`
}

func (p *Person) Create() bool {
	if err := db.Create(p).Error; err != nil {
		return false
	}
	return true
}

func init() {
	db, err = gorm.Open(driver, source)
	if err != nil {
		log.Fatalln(err)
	}
	//db.AutoMigrate(&Person{})
}

func main() {
	p1 := &Person{Name: "王五", Age: 32}
	p1.Create()
	log.Println(p1)
}

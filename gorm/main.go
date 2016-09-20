package main

// http://jinzhu.me/gorm/

import (
        _ "github.com/go-sql-driver/mysql"
        "github.com/jinzhu/gorm"
        "log"
)

type User struct {
        gorm.Model
        Name        string
        Age         uint8
        UserProfile UserProfile                                   // One-To-One relationship (has one - use UserProfile's UserID as foreign key)
        Emails      []Email                                       // One-To-Many relationship (has many - use Email's UserID as foreign key)
        Languages   []Language `gorm:"many2many:user_languages;"` // Many-To-Many relationship, 'user_languages' is join table
}

type UserProfile struct {
        gorm.Model
        // User     User // `UserProfile` belongs to `User`, `UserID` is the foreign key
        // 编译报错： invalid recursive type UserProfile
        // has one 和 belong to 只能定义一个，定义 belong to （One-To-One）
        UserID   uint
        Email    string
        Password string
}

type Email struct {
        gorm.Model
        Email  string
        UserID uint
}

type Language struct {
        gorm.Model
        Name string
}

var driver = "mysql"
var source = `gorm:gorm@tcp(localhost:3306)/gorm?charset=utf8&parseTime=True&loc=Local`

func main() {
        db, err := gorm.Open(driver, source)
        defer db.Close()
        if err != nil {
                log.Fatalln(err)
        }

        //db.AutoMigrate(&User{}, &UserProfile{}, &Email{}, &Language{})

        //db.Create(&Language{Name:"zh-CN"})
        //db.Create(&Language{Name:"zh-TW"})
        //db.Create(&Language{Name:"EN"})
        //db.Create(&Language{Name:"JP"})

        //db.Create(&User{Name:"张三1", Age:34, Languages:[]Language{{Name:"zh-CN"}, {Name:"zh-TW"}}})
        //db.Create(&User{Name:"张三2", Age:35, Languages:[]Language{{Name:"zh-CN"}, {Name:"EN"}}})
        //db.Create(&User{Name:"张三3", Age:36, Languages:[]Language{{Name:"zh-CN"}, {Name:"JP"}}})
        //db.Create(&User{Name:"张三4", Age:37, Languages:[]Language{{Name:"EN"}}})

        var user User
        var languages []Language
        db.First(&user, 1)
        db.Model(&user).Association("Languages").Find(&languages)

        for v, language := range languages {
                log.Println(v, language.ID, language.Name)
        }

        log.Println("end")
}

func main4() {
        db, err := gorm.Open(driver, source)
        defer db.Close()
        if err != nil {
                log.Fatalln(err)
        }

        //db.AutoMigrate(&User{}, &UserProfile{}, &Email{})
        //db.Create(&User{Name:"张三", Age:34})
        //db.Create(&Email{Email:"zhangsan1@localhost.com", UserID:1})
        //db.Create(&Email{Email:"zhangsan2@localhost.com", UserID:1})
        //db.Create(&Email{Email:"zhangsan3@localhost.com", UserID:1})
        //db.Create(&Email{Email:"zhangsan4@localhost.com", UserID:1})

        var user User
        var emails []Email
        db.First(&user, 1)
        db.Model(&user).Related(&emails)

        for v, email := range emails {
                log.Println(v, email.ID, email.Email)
        }

        log.Println("end")
}

func main3() {
        db, err := gorm.Open(driver, source)
        defer db.Close()
        if err != nil {
                log.Fatalln(err)
        }

        db.AutoMigrate(&User{}, &UserProfile{})

        // db.Create(&User{Name:"张三", Age:45})
        // db.Create(&UserProfile{Email:"zhangsan@localhost.com", Password:"123456", UserID:1})

        var user User
        var profile UserProfile
        db.First(&user, 1)
        db.Model(&user).Related(&profile)

        log.Println(user.Name)
        log.Println(profile.Model.ID)

        log.Println("end")
}

func main2() {
        db, err := gorm.Open(driver, source)
        defer db.Close()
        if err != nil {
                log.Fatalln(err)
        }

        db.AutoMigrate(&User{}, &UserProfile{})

        if ok := db.HasTable(&User{}); ok {
                log.Println("ok")
        }

        // 修改字段类型
        db.Model(&User{}).ModifyColumn("Age", "tinyint")

        log.Println("end")
}

func main1() {
        db, err := gorm.Open(driver, source)
        defer db.Close()
        if err != nil {
                log.Fatalln(err)
        }
        log.Println("db open")

        db.AutoMigrate(&User{})

        // create
        // db.Create(&User{Name:"张三", Age:35})

        // read
        var user User
        db.First(&user, 1)

        log.Println(user.CreatedAt.Format("2006-01-02 15:04:05"))
        log.Println(user.Name)

        // update
        db.Model(&user).Update("Age", 40)

        // delete
        var user2 User
        db.First(&user2, 2)
        db.Delete(&user2)

        log.Println("end")
}

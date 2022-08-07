package main

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	var sub = "admin"
	var obj = "book"
	var act = "read"
	e, err := casbin.NewEnforcer("data/model.conf", "data/policy.csv")
	log.Printf("new enforcer %v", err)
	if res, _ := e.Enforce(sub, obj, act); res {
		fmt.Printf("action ok\n")
	} else {
		fmt.Printf("action deny\n")
	}
	UserDB()
	fmt.Printf("done\n")
}

func UserDB() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect database")
		return
	}
	db.AutoMigrate(&Device{})
	db.Create(&Device{
		MAC:    "000EA93D209C",
		Remark: "测试1",
	})
	var dev Device
	db.First(&dev, "mac = ?", "000EA93D209C")
	log.Printf("find device %+v", dev)
}

type Device struct {
	gorm.Model
	MAC    string `json:"MAC"`
	ID     string `json:"ID"`
	Expire string `json:"Expire"`
	Remark string `json:"Remark"`
	Active bool   `json:"Active"`
}

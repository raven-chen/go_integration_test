package main

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/qor/qor"
	"github.com/qor/qor/admin"
)

type User struct {
	gorm.Model
	Name   string
	Gender string
}

var (
	DB gorm.DB
)

func main() {
	Start(9000)
}

func AdminConfig() (mux *http.ServeMux) {
	DB, _ = gorm.Open("sqlite3", "demo.db")
	DB.AutoMigrate(&User{})

	Admin := admin.New(&qor.Config{DB: &DB})
	user := Admin.AddResource(&User{}, &admin.Config{Menu: []string{"User Management"}})
	user.Meta(&admin.Meta{Name: "Gender", Type: "select_one", Collection: []string{"自有设备", "消耗品", "客户设备"}})

	mux = http.NewServeMux()
	Admin.MountTo("/admin", mux)

	return
}

func Start(port int) {
	mux := AdminConfig()
	http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}

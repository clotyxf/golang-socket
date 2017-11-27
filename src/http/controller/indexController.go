package controller

import (
	"database/sql"
	"log"
	"net/http"
	"route"
	"view"

	_ "github.com/go-sql-driver/mysql"
)

var defaults = map[string]bool{
	"/":      true,
	"/index": true,
}

type indexController struct{}

type dataStruct struct {
	Uid  int
	Name string
}

func (self indexController) RegisterRoute() {
	route.HandleFunc("/", self.Index)
}

func (indexController) Index(w http.ResponseWriter, r *http.Request) {
	if _, ok := defaults[r.RequestURI]; !ok {
		http.NotFound(w, r)
		return
	}

	log.Printf("%v", "visit index success")
	_, err := sql.Open("mysql", "cloty:cloty@192.168.137.2:3306/crm_two?charset=utf8")

	if err != nil {
		panic(err)
	}

	data := dataStruct{Uid: 1, Name: "cloty"}
	view.Render(w, r, "index.html", map[string]interface{}{"data": data, "sa": "aa"})
}

package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pintokrysler/when2run/models"
	"github.com/pintokrysler/when2run/server"
)

// IndexController ...
type IndexController struct{}

// NewIndexController creates a new controller instance
func NewIndexController() *IndexController {
	newController := IndexController{}
	return &newController
}

// Dispatch creates a new user instance
func (ic IndexController) Dispatch(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("Index Dispatch!")
	ic.indexHandler(w, r)
}

func (ic IndexController) indexHandler(w http.ResponseWriter, req *http.Request) {
	var myapp = models.App{}
	templateData := models.TplData{
		Title:     "Index",
		TabActive: "index",
	}
	myapp.CurrentView = "index"
	myapp.Data = templateData

	err := server.Server.Tpl.ExecuteTemplate(w, "index.gohtml", myapp)

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error - Resource not found", http.StatusInternalServerError)
	}

}

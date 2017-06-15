package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/pintokrysler/when2run/controllers"
	"github.com/pintokrysler/when2run/models"
	"github.com/pintokrysler/when2run/server"
)

//Myapp ...
var Myapp = models.App{
	UserLogged: false,
}

func init() {
	// Init server
	server.InitServer()
}

func main() {

	router := httprouter.New()

	// Create controllers
	userController := controllers.NewUserController()
	indexController := controllers.NewIndexController()

	router.GET("/", indexController.Dispatch)
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	router.GET("/user/:action", userController.Dispatch)
	router.POST("/user/:action", userController.Dispatch)
	//router.GET("/favicon.ico", http.NotFoundHandler())
	// http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":1700", router)
}

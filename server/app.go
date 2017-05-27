package server

import (
	"database/sql"
	"html/template"
	"log"

	"github.com/gorilla/sessions"
	"github.com/pintokrysler/when2run/models"
)

const (
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "when2run"
)

//Server ..
var Server models.Server

func main() {

}

// InitServer ..
func InitServer() {
	// Init session
	Server.Sess = sessions.NewCookieStore([]byte("something-very-secret"))
	Server.Tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	connectDB()

}

// connectDB ...
func connectDB() {
	databaseConn, err := sql.Open("postgres", "user="+dbUser+" password="+dbPassword+" dbname="+dbName+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	Server.Db = databaseConn
}

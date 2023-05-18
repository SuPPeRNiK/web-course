package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

const (
	port         = ":3000"
	dbDriverName = "mysql"
)

func main() {

	db, err := sql.Open("mysql", "root:1234@tcp(localhost:3306)/blog")
	if err != nil {
		log.Fatal(err)
	}

	dbx := sqlx.NewDb(db, dbDriverName)

	r := mux.NewRouter()
	r.HandleFunc("/home", index(dbx))
	r.HandleFunc("/post/{postID}", post(dbx))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	fmt.Println("Start server " + port)
	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal(err)
	}

}

func openDB() (*sql.DB, error) {
	return sql.Open(dbDriverName, "root:@tcp(localhost:3306)/blog?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true")
}

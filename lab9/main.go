package main

import (
	"log"
	"net/http"
	"text/template"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func database() server {
	database, _ := sql.Open("sqlite3", "parking.db")
	server := server{db: database}
	return server
}

type server struct {
	db *sql.DB
}

func goPage(w http.ResponseWriter, r *http.Request) {
	templ, _ := template.ParseFiles("static/go.html")
	templ.Execute(w, nil)
}

func (s *server) parkingPage(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		err := r.ParseForm()

		if err != nil {
			log.Fatal(err)
		}

		car_num := r.FormValue("car_num")

		period := r.FormValue("period")

		_, err = s.db.Exec("insert into parking(car_num, period) VALUES ($1, $2)", car_num, period)
		if err != nil {
			log.Fatal(err)
		}
		data := map[string]interface{}{"car_num": car_num, "period": period}
		tmpl, _ := template.ParseFiles("static/parking.html")
		tmpl.Execute(w, data)
		return
	}

}

func main() {
	s := database()
	defer s.db.Close()
	fileServer := http.FileServer(http.Dir("./static"))

	http.Handle("/", fileServer)

	http.HandleFunc("/go", goPage)

	http.HandleFunc("/parking", s.parkingPage)

	http.ListenAndServe(":8080", nil)

}

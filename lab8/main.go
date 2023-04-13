package main

import (
	"html/template"
	"log"
	"net/http"
)

func form(w http.ResponseWriter, r *http.Request) {

	tmpl, _ := template.ParseFiles("static/form.html")
	tmpl.Execute(w, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("static/index.html")
	tmpl.Execute(w, nil)
}

func page(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		last := r.FormValue("last")
		data := map[string]interface{}{"name": name, "last": last}
		tmpl, _ := template.ParseFiles("static/page.html")
		tmpl.Execute(w, data)
		return
	}
	tmpl, _ := template.ParseFiles("static/page.html")
	tmpl.Execute(w, nil)
}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))

	http.Handle("/static/", fileServer)
	http.HandleFunc("/", index)
	http.HandleFunc("/form", form)
	http.HandleFunc("/page", page)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}

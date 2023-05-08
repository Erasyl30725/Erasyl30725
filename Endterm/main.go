package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
)

type server struct {
	db *sql.DB
}

type Tour struct {
	Id          int
	Name        string
	Description string
	Price       int
}

type Role struct {
	Id   int
	Name string
}

type User struct {
	Id        int
	FirstName string
	LastName  string
	Phone     string
	Password  string
	Role      string
	Email     string
}

const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "12345"
	dbname   = "besttour"
)

func dbConnect() *server {
	dbconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", dbconn)
	if err != nil {
		log.Fatal(err)
	}
	return &server{db: db}
}

func (s *server) booking(w http.ResponseWriter, r *http.Request) {
	var tours []Tour
	res, _ := s.db.Query("select name, description, price from tours;")
	for res.Next() {
		var tour Tour
		res.Scan(&tour.Name, &tour.Description, &tour.Price)
		tours = append(tours, tour)
	}
	t, _ := template.ParseFiles("static/html/booking.html")
	t.Execute(w, tours)
}

func (s *server) reg(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fn := r.FormValue("fn")
		ln := r.FormValue("ln")
		phone := r.FormValue("phone")
		email := r.FormValue("email")
		pass := r.FormValue("pass")
		_, err := s.db.Exec("insert into users(firstname, lastname, phone, email, password, role) values($1, $2, $3, $4, $5, $6)", fn, ln, phone, email, pass, 1)
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
	}
	t, _ := template.ParseFiles("static/html/signup.html")
	t.Execute(w, nil)
}

func (s *server) auth(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		phone := r.FormValue("phone")
		pass := r.FormValue("pass")
		var passCheck string
		err := s.db.QueryRow("select password from users where phone=$1", phone).Scan(&passCheck)
		if err != nil {
			log.Fatal(err)
		}
		if pass != passCheck {
			fmt.Print("Wrong password")
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	t, _ := template.ParseFiles("static/html/signin.html")
	t.Execute(w, nil)
}

func (s *server) adminAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		phone := r.FormValue("phone")
		pass := r.FormValue("pass")
		var passCheck string
		var role int
		err := s.db.QueryRow("select password, role from users where phone=$1", phone).Scan(&passCheck, &role)
		if err != nil {
			log.Fatal(err)
		}
		if pass != passCheck {
			fmt.Print("Wrong password")
		}
		if role != 2 {
			fmt.Print("You are not admin")
		}
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	t, _ := template.ParseFiles("static/html/adminSignin.html")
	t.Execute(w, nil)
}

func admin(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/html/admin.html")
	t.Execute(w, nil)
}

func (s *server) deleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		_, err := s.db.Exec("delete from users where id=$1", id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Deleted!")
		return
	}
	t, _ := template.ParseFiles("static/html/deleteUser.html")
	t.Execute(w, nil)
}

func (s *server) changeUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fn := r.FormValue("fn")
		ln := r.FormValue("ln")
		phone := r.FormValue("phone")
		role := r.FormValue("role")
		email := r.FormValue("email")
		pass := r.FormValue("pass")
		_, _ = s.db.Exec("update users set firstname=$1, lastname=$2, email=$3, password=$4, role=$5 where phone=$6", fn, ln, email, pass, role, phone)
		http.Redirect(w, r, "/users", http.StatusSeeOther)
	}
	t, _ := template.ParseFiles("static/html/changeUser.html")
	t.Execute(w, nil)
}

func (s *server) users(w http.ResponseWriter, r *http.Request) {
	var users []User
	res, _ := s.db.Query("select * from users;")
	for res.Next() {
		var user User
		res.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Phone, &user.Password, &user.Role, &user.Email)
		users = append(users, user)
	}
	t, _ := template.ParseFiles("static/html/users.html")
	t.Execute(w, users)
}

func (s *server) addTour(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		desc := r.FormValue("desc")
		price := r.FormValue("price")
		_, err := s.db.Exec("insert into tours(name, description, price) values($1, $2, $3)", name, desc, price)
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/booking", http.StatusSeeOther)
	}
	t, _ := template.ParseFiles("static/html/addTour.html")
	t.Execute(w, nil)
}

func (s *server) deleteTour(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		_, err := s.db.Exec("delete from tours where id=$1", id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Deleted!")
		return
	}
	t, _ := template.ParseFiles("static/html/deleteTour.html")
	t.Execute(w, nil)
}

func (s *server) changeTour(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		name := r.FormValue("name")
		desc := r.FormValue("desc")
		price := r.FormValue("price")
		_, err := s.db.Exec("update tours set name=$1, description=$2, price=$3 where id=$4", name, desc, price, id)
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/booking", http.StatusSeeOther)
	}
	t, _ := template.ParseFiles("static/html/changeTour.html")
	t.Execute(w, nil)
}
func (s *server) addRole(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("role")
		_, err := s.db.Exec("insert into roles(role) values($1)", name)
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/roles", http.StatusSeeOther)
	}
	t, _ := template.ParseFiles("static/html/addRole.html")
	t.Execute(w, nil)
}

func (s *server) roles(w http.ResponseWriter, r *http.Request) {
	var roles []Role
	res, _ := s.db.Query("select * from roles;")
	for res.Next() {
		var role Role
		res.Scan(&role.Id, &role.Name)
		roles = append(roles, role)
	}
	t, _ := template.ParseFiles("static/html/roles.html")
	t.Execute(w, roles)
}

func (s *server) deleteRole(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		_, err := s.db.Exec("delete from roles where id=$1", id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Deleted!")
		return
	}
	t, _ := template.ParseFiles("static/html/deleteRole.html")
	t.Execute(w, nil)
}

func main() {
	s := dbConnect()
	defer s.db.Close()

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)

	http.HandleFunc("/booking", s.booking)

	http.HandleFunc("/reg", s.reg)

	http.HandleFunc("/auth", s.auth)

	http.HandleFunc("/adminauth", s.adminAuth)

	http.HandleFunc("/admin", admin)

	http.HandleFunc("/deleteuser", s.deleteUser)

	http.HandleFunc("/changeuser", s.changeUser)

	http.HandleFunc("/users", s.users)

	http.HandleFunc("/addtour", s.addTour)

	http.HandleFunc("/deletetour", s.deleteTour)

	http.HandleFunc("/changetour", s.changeTour)

	http.HandleFunc("/addrole", s.addRole)

	http.HandleFunc("/roles", s.roles)

	http.HandleFunc("/deleterole", s.deleteRole)

	http.ListenAndServe(":8080", nil)
}

package main

import (
	"database/sql"  // Pacote Database SQL para realizar Query
	"log"           // Mostra mensagens no console
	"net/http"      // Gerencia URLs e Servidor Web
	"text/template" // Gerencia templates

	_ "github.com/go-sql-driver/mysql"
)

//Struct utilizada para exibir dados no template

type Person struct {
	Id    int
	Name  string
	Email string
}

// Função dbConn, abre a conexão com o banco de dados

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "crudgo"

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob(("tmpl/*")))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("SELECT * FROM person ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}

	p := Person{}

	res := []Person{}

	for selDB.Next() {

		var id int
		var name, email string

		err = selDB.Scan(&id, &name, &email)

		if err != nil {
			panic(err.Error())
		}

		p.Id = id
		p.Name = name
		p.Email = email

		res = append(res, p)

	}

	tmpl.ExecuteTemplate(w, "Index", res)

	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	pId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT * FROM person WHERE id = ?", pId)

	if err != nil {
		panic(err.Error())
	}

	p := Person{}

	for selDB.Next() {

		var id int
		var name, email string

		err = selDB.Scan(&id, &name, &email)

		if err != nil {
			panic(err.Error())
		}

		p.Id = id
		p.Name = name
		p.Email = email

	}

	tmpl.ExecuteTemplate(w, "Show", p)

	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	pId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT * FROM person WHERE id = ?", pId)

	if err != nil {
		panic(err.Error())
	}

	p := Person{}

	for selDB.Next() {
		var id int
		var name, email string

		err = selDB.Scan(&id, &name, &email)

		if err != nil {
			panic(err.Error())
		}

		p.Id = id
		p.Name = name
		p.Email = email

	}

	tmpl.ExecuteTemplate(w, "Edit", p)

	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	if r.Method == "POST" {

		name := r.FormValue("name")
		email := r.FormValue("email")

		insForm, err := db.Prepare("INSERT INTO person(name, email) VALUES (?,?)")

		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(name, email)

		log.Println("INSERT: Name: " + name + " | E-mail: " + email)
	}

	defer db.Close()

	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	if r.Method == "POST" {

		name := r.FormValue("name")
		email := r.FormValue("email")
		id := r.FormValue("uid")

		updForm, err := db.Prepare("UPDATE person SET name = ?, email = ? WHERE ID = ?")

		if err != nil {
			panic(err.Error())
		}

		updForm.Exec(name, email, id)

		log.Println("UPDATE: Name: " + name + " | E-mail: " + email)
	}

	defer db.Close()

	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	pId := r.URL.Query().Get("id")

	delForm, err := db.Prepare("DELETE FROM person WHERE id = ?")

	if err != nil {
		panic(err.Error())
	}

	delForm.Exec(pId)

	log.Println("DELETE")

	defer db.Close()

	http.Redirect(w, r, "/", 301)
}

func main() {

	log.Println("Server started on: http://localhost:9000")

	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)

	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)

	http.ListenAndServe(":9000", nil)

}

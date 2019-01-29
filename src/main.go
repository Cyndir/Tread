package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

var db *DB

func main() {

	db, err := sql.Open("sqlite3", "./db.sql")
	checkErr(err)
	defer db.Close()
	var router = mux.NewRouter()
	router.HandleFunc("/healthcheck", healthCheck).Methods("GET")
	router.HandleFunc("/getall/{user}", returnAll).Methods("GET")
	router.HandleFunc("/delete", handleDelete).Methods("DELETE")
	router.HandleFunc("/add", handleAdd).methods("PUT")

	fmt.Println("Running server!")
	//Set up graceful shutdown
	errs := make(chan error, 2)
	server := &http.Server{Addr: ":8080", Handler: handler}

	go func() {
		errs <- server.ListenAndServe()
	}()

	go func() {
		// Setting up signal capturing
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt)
		errs <- fmt.Errorf("%s", <-c)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	// Waiting for SIGINT (pkill -2)
	<-errs
}
func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Still alive!")
}

func returnAll(w http.ResponseWriter, r *http.Request) {
	var uid int
	var link string
	var links []string
	uidSelect, err := db.Prepare("SELECT uid FROM users where uname = ?")
	checkErr(err)
	linkSelect, err := db.Prepare("SELECT link FROM links WHERE uid = ?")
	checkErr(err)

	uName := mux.Vars(r)["user"]
	rows, err := uidSelect.Query(uName)
	checkErr(err)
	rows.Next()
	err = rows.Scan(&uid)
	checkErr(err)

	rows, err = linkSelect.Query(uid)
	checkErr(err)
	for rows.Next() {
		rows.Scan(&link)
		append(links, link)
	}
	json.NewEncoder(w).Encode(links)

}
func handleDelete(w http.ResponseWriter, r *http.Request) {

	var uid int

	vars := r.URL.Query()
	uname := vars.Get("name")
	link := vars.Get("link")
	//would need to check auth here
	uidSelect, err := db.Prepare("SELECT uid FROM users where uname = ?")
	checkErr(err)
	linkDelete, err := db.Prepare("DELETE FROM links WHERE uid = ? AND link = ?")
	checkErr(err)

	rows, err := uidSelect.Query(uName)
	checkErr(err)
	rows.Next()
	err = rows.Scan(&uid)
	checkErr(err)

	res, err := linkDelete.Exec(uid, link)
	checkErr(err)
	msg := fmt.Sprintf("Rows deleted: %d", res.RowsAffected())
	json.NewEncoder(w).Encode(map[string]string{"message": msg})

}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	var uid int

	vars := r.URL.Query()
	uname := vars.Get("name")
	link := vars.Get("link")

	uidSelect, err := db.Prepare("SELECT uid FROM users where uname = ?")
	checkErr(err)
	linkAdd, err := db.Prepare("INSERT INTO links (uid, link) VALUES (?, ?)")
	checkErr(err)

	rows, err := uidSelect.Query(uName)
	checkErr(err)
	rows.Next()
	err = rows.Scan(&uid)
	checkErr(err)

	res, err := linkAdd.Exec(uid, link)
	checkErr(err)
	msg := fmt.Sprintf("Rows added: %d", res.RowsAffected())
	json.NewEncoder(w).Encode(map[string]string{"message": msg})

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

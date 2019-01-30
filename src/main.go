package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"os"
	"os/signal"
)

var db *sql.DB

func main() {

	db, err := sql.Open("sqlite3", "./db.sql")
	checkErr(err)
	defer db.Close()
	var router = mux.NewRouter()
	router.HandleFunc("/healthcheck", healthCheck).Methods("GET")
	router.HandleFunc("/getall/{user}", getAll).Methods("GET")
	router.HandleFunc("/delete", handleDelete).Methods("DELETE")
	router.HandleFunc("/add", handleAdd).Methods("PUT")

	fmt.Println("Running server!")
	//Set up graceful shutdown
	errs := make(chan error, 2)
	server := &http.Server{Addr: ":8080", Handler: router}

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
	<-errs
}
func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Still alive!")
}

func getAll(w http.ResponseWriter, r *http.Request) {
	var link string
	var links []string
	uname := mux.Vars(r)["user"]
	if uname == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	linkSelect, err := db.Prepare("SELECT link FROM links WHERE uid = ?")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	uid, err := getUid(uname)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rows, err := linkSelect.Query(uid)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		rows.Scan(&link)
		links = append(links, link)
	}
	json.NewEncoder(w).Encode(links)

}
func handleDelete(w http.ResponseWriter, r *http.Request) {

	vars := r.URL.Query()
	uname := vars.Get("name")
	link := vars.Get("link")
	if link == "" || uname == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	//would need to check auth here
	linkDelete, err := db.Prepare("DELETE FROM links WHERE uid = ? AND link = ?")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	uid, err := getUid(uname)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	res, err := linkDelete.Exec(uid, link)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	numRows, err := res.RowsAffected()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	msg := fmt.Sprintf("Rows deleted: %d", numRows)
	json.NewEncoder(w).Encode(map[string]string{"message": msg})

}

func handleAdd(w http.ResponseWriter, r *http.Request) {

	vars := r.URL.Query()
	uname := vars.Get("name")
	link := vars.Get("link")
	if link == "" || uname == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	linkAdd, err := db.Prepare("INSERT INTO links (uid, link) VALUES (?, ?)")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	uid, err := getUid(uname)
	res, err := linkAdd.Exec(uid, link)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	numRows, err := res.RowsAffected()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	msg := fmt.Sprintf("Rows added: %d", numRows)
	json.NewEncoder(w).Encode(map[string]string{"message": msg})

}
func getUid(uName string) (uid int, err error) {
	uidSelect, err := db.Prepare("SELECT uid FROM users where uname = ?")
	if err != nil {
		return
	}
	row := uidSelect.QueryRow(uName)
	err = row.Scan(&uid)
	if err != nil {
		return
	}

	return

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

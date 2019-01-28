package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	var router = mux.NewRouter()
	router.HandleFunc("/healthcheck", healthCheck).Methods("GET")
	router.HandleFunc("/getall/{user}", returnAll).Methods("GET")
	router.HandleFunc("/delete", handleDelete).Methods("DELETE")
	router.HandleFunc("/add", handleAdd).methods("PUT")

	fmt.Println("Running server!")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Still alive!")
}

func returnAll(w http.ResponseWriter, r *http.Request) {

}
func handleDelete(w http.ResponseWriter, r *http.Request) {

}

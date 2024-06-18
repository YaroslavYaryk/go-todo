package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UP and running")
}

func main() {

	var router = mux.NewRouter()
	const port = ":8080"

	router.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(port, router))

}

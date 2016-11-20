package main

import (
	"fmt"
	"net/http"

	"log"

	"github.com/gorilla/mux"
)

var counter int

func init() {
	counter = 0
}

func startServer(port int) {
	r := mux.NewRouter()
	r.HandleFunc("/greetings", handleGreetings)

	bind := fmt.Sprintf("localhost:%d", port)
	http.ListenAndServe(bind, r)
}

func handleGreetings(w http.ResponseWriter, r *http.Request) {
	counter++
	str := fmt.Sprintf("greetings:%d ", counter)

	log.Print(str)
	w.Write([]byte(str))
}

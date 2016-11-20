package main

import (
	"fmt"
	"net/http"

	"log"

	"github.com/gorilla/mux"
)

var counter int
var httpPort int

func init() {
	counter = 0
}

func startServer(port int) {
	httpPort = port

	r := mux.NewRouter()
	r.HandleFunc("/greetings", handleGreetings)

	bind := fmt.Sprintf("localhost:%d", httpPort)

	log.Printf("Binding server at %s", bind)
	http.ListenAndServe(bind, r)
}

func handleGreetings(w http.ResponseWriter, r *http.Request) {
	counter++
	str := fmt.Sprintf("greetings[%d]:%d\n", httpPort, counter)

	log.Print(str)
	w.Write([]byte(str))
}

package main

import (
	"fmt"
	"net/http"
	"net/url"

	"log"

	"os"

	"net/http/httputil"
)

var counter int
var httpPort int

func init() {
	counter = 0
}

func startServer(port int) {
	httpPort = port

	if flags.Plugin {
		http.HandleFunc("/greetings", handleGreetings)
	} else {
		log.Print("setting up reverse proxy")

		pluginPort := flags.Port + 1
		proxyUrl, _ := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", pluginPort))

		reverseProxy := httputil.NewSingleHostReverseProxy(proxyUrl)
		http.HandleFunc("/greetings", reverseProxy.ServeHTTP)
	}

	bind := fmt.Sprintf(":%d", httpPort)

	log.Printf("Binding server at %s", bind)
	http.ListenAndServe(bind, nil)
}

func handleGreetings(w http.ResponseWriter, r *http.Request) {
	counter++
	str := fmt.Sprintf("greetings[%d]:%d\n", httpPort, counter)

	//log.Print(str)
	fmt.Fprintf(os.Stdout, "%s", str)
	fmt.Fprintf(os.Stderr, "%s", str)

	w.Write([]byte(str))
}

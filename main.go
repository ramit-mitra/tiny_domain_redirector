package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const redirectFrom = "contact.ramit.io"
const redirectTo = "https://ramit.io/contact"

func main() {
	// create a file to store the logs
	file, err := os.Create("app.log")
	if err != nil {
		fmt.Println(err)
	}

	// set the output of the logger to the file
	log.SetOutput(file)

	// set app port
	port := flag.Int("port", 9990, "Port to listen on")
	flag.Parse()

	log.Println("--- Starting redirector")
	log.Println("--- Running on port: " + fmt.Sprint(*port))

	// default handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// let's set some headers
		w.Header().Set("cache-control", "public, max-age=3600")
		w.Header().Set("service", "ramit/tiny_domain_redirector")
		w.Header().Set("project-url", "https://github.com/ramit-mitra/tiny_domain_redirector")

		if strings.Contains(r.Host, redirectFrom) {
			log.Println("✅ Redirecting to " + fmt.Sprint(redirectTo))
			http.Redirect(w, r, redirectTo, http.StatusFound)
		} else {
			log.Println("❌ Not redirecting, request from host: " + fmt.Sprint(r.Host))

			// return a HTTP status code of 404
			w.WriteHeader(http.StatusNotFound)
		}
	})

	// start app server
	http.ListenAndServe("127.0.0.1:"+fmt.Sprint(*port), nil)
}

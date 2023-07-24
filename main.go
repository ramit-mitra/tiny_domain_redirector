package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// define redirection map to store redirection details
	redirectionMap := map[string]string{
		"contact.ramit.io": "https://ramit.io/contact",
		"ramitmitra.com":   "https://ramit.io",
		"ramit.org":        "https://ramit.io",
		"ramit.dev":        "https://ramit.io",
		"ramit.app":        "https://ramit.io",
	}

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

	log.Println("--- 🚀 Starting redirector ")
	log.Println("--- 👾 Running on port: " + fmt.Sprint(*port))

	// default handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// let's set some headers
		w.Header().Set("cache-control", "private, no-cache")
		w.Header().Set("service", "ramit/tiny_domain_redirector")
		w.Header().Set("project-url", "https://github.com/ramit-mitra/tiny_domain_redirector")

		if target, ok := redirectionMap[r.Host]; ok {
			log.Printf("✅ Redirecting from %s to %s\n", fmt.Sprint(r.Host), fmt.Sprint(target))
			http.Redirect(w, r, target, http.StatusFound)
		} else {
			log.Printf("❌ Not redirecting, request from host: %s\n", fmt.Sprint(r.Host))
			// return a HTTP status code of 404
			w.WriteHeader(http.StatusNotFound)
		}
	})

	// start app server
	http.ListenAndServe("127.0.0.1:"+fmt.Sprint(*port), nil)
}

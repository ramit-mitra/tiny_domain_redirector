package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"
)

const Version = "1.1.1"

type Redirect struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
	StatusCode  int    `yaml:"status_code,omitempty"`
}

type Config struct {
	Redirects []Redirect `yaml:"redirects"`
}

func main() {
	// create a file to store the logs
	file, err := os.Create("app.log")
	if err != nil {
		fmt.Println(err)
	}

	// set the output of the logger to the file
	log.SetOutput(file)

	// create a file to store metrics
	metricsFile, err := os.OpenFile("metrics.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer metricsFile.Close() // Close the metrics file when main exits

	// create a new logger for metrics
	metricsLogger := log.New(metricsFile, "", 0) // No prefix, no flags (timestamp, etc.)

	// set app port
	port := flag.Int("port", 9990, "Port to listen on")
	flag.Parse()

	// read the redirects.yaml file
	yamlFile, err := os.ReadFile("redirects.yaml")
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Error reading redirects.yaml")
		fmt.Fprintf(os.Stderr, "😭 %v\n", err)
		os.Exit(101)
	}

	// unmarshal the yaml file into a Config struct
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Error unmarshalling redirects.yaml")
		fmt.Fprintf(os.Stderr, "😭 %v\n", err)
		os.Exit(110)
	}

	// create redirection map to store redirection details
	redirectionMap := make(map[string]Redirect)
	if len(config.Redirects) == 0 {
		fmt.Fprintln(os.Stderr, "❌ No redirects found in redirects.yaml. Exiting.")
		os.Exit(111)
	}
	for _, r := range config.Redirects {
		redirectionMap[r.Source] = r
	}

	log.Println("--- 🚀 Starting redirector (version: " + fmt.Sprint(Version) + ")")
	log.Println("--- 👾 Running on port: " + fmt.Sprint(*port))

	// default handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		metricsLogger.Println("http_requests_total 1") // Log every request

		// let's set some headers
		w.Header().Set("cache-control", "private, no-cache")
		w.Header().Set("service", "ramit/tiny_domain_redirector")
		w.Header().Set("project-url", "https://github.com/ramit-mitra/tiny_domain_redirector")

		if redirect, ok := redirectionMap[r.Host]; ok {
			// default to 302 if status code is not specified
			statusCode := redirect.StatusCode
			if statusCode == 0 {
				statusCode = http.StatusFound
			}
			log.Printf("✅ Redirecting from %s to %s\n", fmt.Sprint(r.Host), fmt.Sprint(redirect.Destination))
			metricsLogger.Printf("http_redirects_total{host=\"%s\",destination=\"%s\",status_code=\"%d\"} 1\n", r.Host, redirect.Destination, statusCode) // Log successful redirect
			http.Redirect(w, r, redirect.Destination, statusCode)
		} else {
			log.Printf("❌ Not redirecting, request from host: %s\n", fmt.Sprint(r.Host))
			// return a HTTP status code of 404
			w.WriteHeader(http.StatusNotFound)
		}
	})

	// healthz endpoint
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		metricsLogger.Println("http_requests_total{path=\"/healthz\"} 1") // Log healthz request
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	// create a channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// create a server
	server := &http.Server{
		Addr:    "127.0.0.1:" + fmt.Sprint(*port),
		Handler: nil, // use default http.ServeMux
	}

	// start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	log.Println("--- ✅ Server is ready to handle requests")

	// wait for a signal
	<-quit
	log.Println("---  shutting down server...")

	// create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server shutdown failed: %v", err)
	}

	log.Println("--- ✅ Server gracefully stopped")
}

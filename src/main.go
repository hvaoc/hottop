package main

import (
	"fmt"
	"net/http"
	"flag"
	"log"
	"strconv"
	"os"
	"io/ioutil"
	"strings"
	"os/signal"
	"syscall"
)

var WORKING_DIR  = ""

func getFile(rd string, filePath string, extension string) ([]byte, error)  {

	if strings.HasSuffix(filePath, "." + extension) {
		filePath = rd + filePath
	} else {
		filePath = rd + filePath + "." + extension
	}

	// Read a file & return error when failed to read
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func handleRequest(w http.ResponseWriter, r *http.Request)  {

	// Parse Form
	// r.ParseForm()

	// Clean request URL path
	requestPath := r.URL.Path
	last := len(requestPath) - 1
	if requestPath[last] == '/' {
		requestPath = requestPath[:last]
	}


	// Try HTML
	raw, err := getFile(WORKING_DIR, requestPath, "html" )
	if err == nil {
		log.Printf("| 200 | GET %s", requestPath)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, string(raw[:]))
	} else {
		// Try JSON
		raw, err := getFile(WORKING_DIR, requestPath, "json" )
		if err == nil {
			log.Printf("| 200 | GET %s", requestPath)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, string(raw[:]))
		} else {
			// Try XML
			raw, err := getFile(WORKING_DIR, requestPath, "xml" )
			if err == nil {
				log.Printf("| 200 | GET %s", requestPath)
				w.Header().Set("Content-Type", "application/xml")
				fmt.Fprintf(w, string(raw[:]))
			} else {
				log.Printf("| 404 | GET %s", requestPath)
				http.NotFound(w, r)
			}
		}
	}

}

func startServer(port int)  {

	// Assign IP & port on which the server will listen on
	addr := ":" + strconv.Itoa(port)

	// Assign handler to process incoming requests
	http.HandleFunc("/", handleRequest) // set router

	// Try to start listening to requests
	err := http.ListenAndServe(addr, nil) // set listen port
	if err != nil { // Unable to start the server
		log.Fatal("Failed starting server: ", err)
	}
}

func setGlobalVariables()  {
	// Get current working directory
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	WORKING_DIR = pwd
}

func cleanup()  {
	fmt.Println("Shutting down server - [DONE]")
}

func main() {
	// Set Global Variables
	setGlobalVariables()

	// Manage Interrupts
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	// Print banner
	version := "0.0.1"
	fmt.Println("Hottop : HTTP Server - " + version)
	port := flag.Int("port", 8080, "Port Number")

	// Parse command line flags
	flag.Parse()

	fmt.Println("Listening on port:", *port)
	fmt.Println("Root directory:", WORKING_DIR)

	// Start HTTP server and listen to incoming requests
	startServer(*port)

}


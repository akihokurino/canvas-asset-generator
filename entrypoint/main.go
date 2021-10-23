package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.DefaultServeMux

	mux.HandleFunc("/", apiHandler)
	mux.HandleFunc("/video_split", taskHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("running server on port: %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("failed running server, err=%+v", err)
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "Hello, World!")
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ReadAll: %v", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	output := fmt.Sprintf("Completed task: payload(%s)", string(body))
	log.Println(output)

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintln(w, output)
}

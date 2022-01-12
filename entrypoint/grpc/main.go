package main

import (
	"canvas-server/di"
	"canvas-server/entrypoint"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	if os.Getenv("IS_LOCAL") == "true" {
		entrypoint.MustLoadLocalEnv("/app/app.yaml")
	}

	mux := http.DefaultServeMux

	di.ResolveGRPC()(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("running server on port: %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("failed running server, err=%+v", err)
	}
}

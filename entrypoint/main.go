package main

import (
	"canvas-server/di"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

func main() {
	if os.Getenv("IS_LOCAL") == "true" {
		MustLoadLocalEnv("/app/app.yaml")
	}

	mux := http.DefaultServeMux

	di.ResolveSubscriber()(mux)
	di.ResolveBatch()(mux)
	di.ResolveGraphQL()(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("running server on port: %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("failed running server, err=%+v", err)
	}
}

type env struct {
	Variables map[string]string `yaml:"env_variables"`
}

func MustLoadLocalEnv(path string) {
	if os.Getenv("IS_LOCAL") != "true" {
		return
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	e := env{}
	if err := yaml.Unmarshal(buf, &e); err != nil {
		panic(err)
	}

	for k, v := range e.Variables {
		if err := os.Setenv(k, v); err != nil {
			panic(err)
		}
	}
}

package main

import (
	"flag"
	"net/http"

	"strconv"

	"spark-cluster/dashboard/backend/handler"

	log "github.com/sirupsen/logrus"
)

const (
	DefaultFrontendDir string = "dashboard/frontend"
	DefaultListenPort  int    = 8081
)

var (
	frontendDir string
	listenPort  int
)

func main() {
	flag.StringVar(&frontendDir, "frontend-dir", DefaultFrontendDir, `directory of the dashboard frontend`)
	flag.IntVar(&listenPort, "port", DefaultListenPort, `port this server listen to`)
	flag.Parse()

	h, err := handler.NewAPIHandler(frontendDir)
	if err != nil {
		log.Fatalf("Failed to create route handler: %v", err)
	}
	router := handler.NewRouter(h)

	http.Handle("/", router)

	p := ":" + strconv.Itoa(listenPort)
	log.Infof("Spark-cluster backend server listens on: %v", p)
	if err = http.ListenAndServe(p, nil); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

package main

import (
	"log"
	"net/http"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/health"
)

func healthCheck() http.HandlerFunc {
	return health.Check
}

func main() {
	router := http.NewServeMux()
	router.Handle("/", http.FileServer(http.Dir("public")))
	router.HandleFunc("/api/health/web", healthCheck())

	log.Fatal(http.ListenAndServe("0.0.0.0:3000", router))
}

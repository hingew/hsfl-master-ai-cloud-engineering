package main

import (
	"log"
	"net/http"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/request-service/coalecing"
)

func main() {
	coalecer := coalecing.NewWebServiceCoalecer()

	if err := http.ListenAndServe(":3000", coalecer); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}

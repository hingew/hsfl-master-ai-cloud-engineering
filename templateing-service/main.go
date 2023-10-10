package main

import (
	"log"
	"net/http"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/router"
)

func main() {
	handler := router.New()
	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}

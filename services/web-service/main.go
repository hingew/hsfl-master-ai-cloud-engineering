package main

import (
	"log"
	"net/http"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/web-service/api/router"
)

const DIR = "./public/"

func main() {
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", &router.Router{}))
}

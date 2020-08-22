package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/svasandani/terrabaq/internal/api"
	"github.com/svasandani/terrabaq/internal/db"
)

func main() {
	pukkalink := flag.String("pukkaLink", "https://pukka.terraling.com/", "Link for email redirect")

	port := flag.String("port", "3000", "Port to serve Terrabaq")

	flag.Parse()

	sessions := make(map[string]db.User)

	api.SetupAPI(pukkalink, sessions)

	http.HandleFunc("/", api.SessionHandler)
	http.HandleFunc("/enqueue", api.EnqueueHandler)

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

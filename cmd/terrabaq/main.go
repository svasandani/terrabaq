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
	usrtouuid := make(map[string]string)

	api.SetupAPI(pukkalink, sessions, usrtouuid)

	http.HandleFunc("/new_session", api.Middleware(api.SessionHandler))
	http.HandleFunc("/enqueue", api.Middleware(api.EnqueueHandler))

	http.HandleFunc("/update_roles", api.Middleware(api.UpdateHandler))

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"example.com/rest-api-raw/admin"
	"example.com/rest-api-raw/coaster"
)

func main() {

	// start the service with env var
	// ADMIN_PASSWORD=secret go run main.go

	// init Admin
	ad := admin.NewAdminPortal()
	http.HandleFunc("/admin", ad.Handler)

	// init
	var cstrStore coaster.Store
	// set a dummy coaster
	now := fmt.Sprintf("%d", time.Now().UnixNano())
	cstrStore.Store = map[string]coaster.Coaster{
		now: coaster.Coaster{
			ID:           now,
			Name:         "Fury 325",
			Manufacturer: "B+M",
			InPark:       "Carowinds",
			Height:       99,
		},
	}
	http.Handle("/coaster", &cstrStore)
	http.Handle("/coaster/", &cstrStore)

	// Use the default Server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

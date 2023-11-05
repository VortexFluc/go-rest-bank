package main

import (
	"log"
)

func main() {
	store, storeErr := NewPostgresStore()
	if storeErr != nil {
		log.Fatal(storeErr)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":8008", store)
	runErr := server.Run()
	if runErr != nil {
		log.Fatal(runErr)
	}

}

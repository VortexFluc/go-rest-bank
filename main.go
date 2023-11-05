package main

import (
	"fmt"
	"log"
)

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", store)
	//server := NewAPIServer(":8008", store)
	//server.Run()

}

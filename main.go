package main

import (
	"flag"
	"fmt"
	"log"
)

func seedAccount(store Storage, fname, lname, pw string) *Account {
	acc, err := NewAccount(fname, lname, pw)
	if err != nil {
		log.Fatal(err)
	}

	if err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	fmt.Println("new account -> ", acc.Number)
	return acc
}

func seedAccounts(s Storage) {
	seedAccount(s, "andrey", "v", "123456")
}

func main() {
	seed := flag.Bool("seed", false, "seed the db")
	flag.Parse()

	store, storeErr := NewPostgresStore()
	if storeErr != nil {
		log.Fatal(storeErr)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Printf("Seeding the DB\n")
		// seed stuff
		seedAccounts(store)
	}

	server := NewAPIServer(":8008", store)
	runErr := server.Run()
	if runErr != nil {
		log.Fatal(runErr)
	}

}

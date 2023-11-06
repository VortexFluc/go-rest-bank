package main

import (
	"flag"
	"fmt"
	"github.com/labstack/gommon/color"
	"github.com/vortexfluc/gobank/internal/gobank/account"
	"github.com/vortexfluc/gobank/internal/gobank/api"
	"github.com/vortexfluc/gobank/internal/gobank/storage"
	"log"
)

const banner = `

  ______             _______                       __       
 /      \           /       \                     /  |      
/$$$$$$  |  ______  $$$$$$$  |  ______   _______  $$ |   __ 
$$ | _$$/  /      \ $$ |__$$ | /      \ /       \ $$ |  /  |
$$ |/    |/$$$$$$  |$$    $$<  $$$$$$  |$$$$$$$  |$$ |_/$$/ 
$$ |$$$$ |$$ |  $$ |$$$$$$$  | /    $$ |$$ |  $$ |$$   $$<  
$$ \__$$ |$$ \__$$ |$$ |__$$ |/$$$$$$$ |$$ |  $$ |$$$$$$  \ 
$$    $$/ $$    $$/ $$    $$/ $$    $$ |$$ |  $$ |$$ | $$  |
 $$$$$$/   $$$$$$/  $$$$$$$/   $$$$$$$/ $$/   $$/ $$/   $$/ 
                                                            
[Designed by %s]                                                            
                                                            

`

func seedAccount(store storage.Storage, fname, lname, pw string) *account.Account {
	acc, err := account.NewAccount(fname, lname, pw)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	fmt.Println("new account -> ", acc.Number)
	return acc
}

func seedAccounts(s storage.Storage) {
	seedAccount(s, "andrey", "v", "123456")
}

func main() {
	seed := flag.Bool("seed", false, "seed the db")
	flag.Parse()

	store, storeErr := storage.NewPostgresStore()
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

	colorer := color.New()
	colorer.Printf(banner, colorer.Blue("VortexFluc"))

	server := api.NewAPIServer(":8008", store)
	runErr := server.Run()
	if runErr != nil {
		log.Fatal(runErr)
	}

}

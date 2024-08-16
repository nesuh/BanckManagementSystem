package main

import (
	"fmt"
	"log"
)

func main() {
	store, err := NewPostgressStore()
	if err != nil {
		log.Fatal(err)
	}
	defer store.db.Close() // Ensure the database connection is closed when done

	if err := store.createAccountTable(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIserver(":3000", store)
	server.Run()

	fmt.Printf("%+v\n", store)
}

package main

import "github.com/fomentia/ascendhosting/app"

func main() {
	var err error
	db, err = database.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	app.Serve(db)
}

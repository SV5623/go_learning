package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// This is a placeholder for the actual migration code.
	// In a real implementation, you would set up your database connection,
	// define your migration`s, and execute them here.
	if len (os.Args) < 2 {
		log.Fatal("Please provide a migration direction: `up` or `down`.")
	}
	direction := os.Args[1]
	db, err := sql.Open(
		"postgres",
		"postgres://postgres:postgres@localhost:5432/event_app?sslmode=disable",
	)
	if err != nil {
		panic(err)
	}
	defer db.Close()


	instance, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	
	cwd, _ := os.Getwd()
	log.Println("Current dir:", cwd)

	fSrc, err := (&file.File{}).Open("./migrations")
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithInstance("file", fSrc, "postgres", instance)
	if err != nil {
		log.Fatal(err)
	}
	switch direction {	
		case "up":
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				log.Fatal(err)
			}
			log.Println("Migrations applied successfully.")
		case "down":
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				log.Fatal(err)
			}
			log.Println("Migrations rolled back successfully.")
		default:
			log.Fatal("Invalid migration direction. Use `up` or `down`.")	
	}
}

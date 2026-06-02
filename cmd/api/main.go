package main

import (
	"database/sql"
	"log"
	"rest_app_in_gin/internal/database"
	"rest_app_in_gin/internal/env"
	_ "github.com/lib/pq"
	_ "github.com/joho/godotenv/autoload"
)


type application struct {
	port int
	jwtSecret string
	models database.Models
}


func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/event_app?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	models := database.NewModels(db)
	app := &application{
		port: env.GetEnvInt("PORT", 5623),
		jwtSecret: env.GetEnvString("JWT_SECRET", "defaultsecret"),
		models: models,
	}

	log.Printf("Starting server on port %d", app.port)
	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}

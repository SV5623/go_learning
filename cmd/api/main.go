package main

import (
	"database/sql"
	"log"
	"os"
	"rest_app_in_gin/internal/database"
	"strconv"

	//	"rest_app_in_gin/internal/env"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)


type application struct {
	port int
	jwtSecret string
	models database.Models
}


func main() {
	db, err := sql.Open(
	"postgres",
	os.Getenv("DATABASE_URL"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}

	models := database.NewModels(db)
	app := &application{
		port: port,
		jwtSecret: os.Getenv("JWT_SECRET"),
		models: models,
	}
	log.Println("PORT =", os.Getenv("PORT"))
	log.Println("DATABASE_URL =", os.Getenv("DATABASE_URL"))
	log.Printf("Starting server on port %d", app.port)
	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}

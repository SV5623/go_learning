package main

import (
	"api_project/internal/env"
	"context"
	"log"
	"log/slog"
	"os"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	ctx := context.Background()

	cfg := &Config{
		addr : ":2000",
		db: dbConfig{		
			dsn: env.GetString("GOOSE_DBSTRING", ""),
		},	
	}	
	if cfg.db.dsn == "" {
		log.Fatal("GOOSE_DBSTRING is not set")
	}
	//database
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		log.Printf("Помилка підключення до бази даних: %s", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	api := aplication{
	config: *cfg,
	}
	
	h := api.mount() //повертає хендлер який ми будемо запускати
	if err := api.run(h); err != nil { // запускає сервер і слухає на вказаному адресі
		// slog.Error("Сервер не вдалося запустити:, помилка %s", err) 
		log.Printf("Сервер не вдалося запустити:, помилка %s", err) 
		// log.Println("Сервер не вдалося запустити:, помилка %s", err)
		os.Exit(1)
	}
	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// або мож написати інакше
	// if err := api.run(api.mount()); err != nil { // запускає сервер і слухає на вказаному адресі
	// 	log.Fatal(err)
	// }
}

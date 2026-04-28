package main

import (
	"log"
	"log/slog"
	"os"
)

func main() {
	cfg := &Config{
		addr : ":2000",
		db: dbConfig{},
	}

	api := aplication{
	config: *cfg,
	}

	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	slog.SetDefault(logger)

	h := api.mount() //повертає хендлер який ми будемо запускати
	if err := api.run(h); err != nil { // запускає сервер і слухає на вказаному адресі
		// slog.Error("Сервер не вдалося запустити:, помилка %s", err) 
		log.Printf("Сервер не вдалося запустити:, помилка %s", err) 
		// log.Println("Сервер не вдалося запустити:, помилка %s", err)
		os.Exit(1)
	}

	// або мож написати інакше
	// if err := api.run(api.mount()); err != nil { // запускає сервер і слухає на вказаному адресі
	// 	log.Fatal(err)
	// }
}

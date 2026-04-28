package main

import (
	"api_project/internal/products"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


type aplication struct {
	config Config
	// logger
	// db driver
}

func (app *aplication) mount() http.Handler {
	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi\n"))
	})

	// Products routes з файлу handlers.go
	productsService := products.NewService() // тут ми створюємо сервіс який реалізує інтерфейс Service, але поки що ми передаємо nil, бо ще не реалізували сервіс
	productsHandler := products.NewHandler(productsService) // тут треба передати сервіс який буде реалізовувати інтерфейс Service, але поки що ми передаємо nil, бо ще не реалізували сервіс
	r.Get("/products", productsHandler.ListProductsHandler) // тут ми вказуємо що при запиті на /products буде викликатися метод ListProductsHandler який знаходиться в productsHandler
	

	
	return r
}

func (app *aplication) run(h http.Handler) error { //h це хендлер який ми будемо запускати
	srv := &http.Server{ // srv це сервер який ми будемо запускати
		Addr:         app.config.addr,
		Handler:      h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: time.Minute * 1, // якщо написати просто time.Minute то буде 1 хвилина, якщо написати time.Minute * 1 то буде теж саме, але можна написати time.Minute * 5 і буде 5 хвилин
	}
	log.Printf("Сервер запустився на адресі %s", app.config.addr)
	
	return srv.ListenAndServe() // ListenAndServe запускає сервер і слухає на вказаному адресі, якщо виникає помилка, вона повертається, якщо ж ні то сервер працює і слухає на вказаному адресі
}


type Config struct {
	addr string
	db dbConfig
}

type dbConfig struct {
	dsn string //domain string user== password== dbname== 
}
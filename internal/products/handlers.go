package products

import (
	//"api_project/internal/products"
	// "encoding/json"
	"api_project/internal/json"
	"log"
	"net/http"
)

type handler struct {
	service Service
	// logger
	// db driver
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for getting a product
	products ,err := h.service.ListProducts(r.Context()) // тут ми викликаємо метод ListProducts який знаходиться в нашому сервісі і передаємо йому контекст з запиту, якщо виникає помилка, ми повертаємо її клієнту
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError) // тут ми використовуємо функцію http.Error яка знаходиться в пакеті net/http і передаємо їй відповідь яку ми хочемо повернути клієнту, статус відповіді і дані які ми хочемо повернути клієнту
		return
	}

	json.WriteJSON(w, http.StatusOK, products) // тут ми використовуємо нашу функцію WriteJSON яка знаходиться в internal/json/json.go і передаємо їй відповідь яку ми хочемо повернути клієнту, статус відповіді і дані які ми хочемо повернути клієнту

}

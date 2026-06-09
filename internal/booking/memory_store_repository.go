package booking

// MemoryStore - проста "база даних" в оперативній пам'яті.
// bookings зберігає всі бронювання.
//
// Ключ (string)   -> SeatID (ідентифікатор місця)
// Значення        -> Booking (інформація про бронювання)
type MemoryStore struct {
	bookings map[string]Booking
}

// NewMemoryStore створює новий MemoryStore.
//
// Це конструктор.
// Він створює порожню map, куди потім будуть записуватись бронювання.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		bookings: map[string]Booking{},
	}
}

// Book додає нове бронювання.
//
// Працює так:
// 1. Перевіряє чи є вже запис з таким SeatID.
// 2. Якщо є - повертає помилку.
// 3. Якщо немає - записує бронювання в map.
func (s *MemoryStore) Book(b Booking) error {

	// Шукаємо місце в map.
	//
	// _       -> саме значення нам не потрібне.
	// exists  -> true якщо ключ знайдений.
	if _, exists := s.bookings[b.SeatID]; exists {
		return ErrSeatsAlreadyBooked
	}

	// Додаємо бронювання.
	//
	// Ключом буде SeatID.
	// Значенням буде весь об'єкт Booking.
	s.bookings[b.SeatID] = b

	return nil
}

// ListBookings повертає всі бронювання для конкретного фільму.
//
// movieID - ідентифікатор фільму,
// для якого потрібно знайти бронювання.
func (s *MemoryStore) ListBookings(movieID string) []Booking {

	// Порожній зріз для результату.
	var result []Booking

	// Перебираємо всі бронювання в map.
	for _, b := range s.bookings {

		// Якщо бронювання належить потрібному фільму,
		// додаємо його в результат.
		if b.MovieID == movieID {
			result = append(result, b)
		}
	}

	// Повертаємо список знайдених бронювань.
	return result
}
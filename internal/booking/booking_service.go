package booking

// Service - бізнес-логіка застосунку.
//
// Сам Service не знає,
// де саме зберігаються дані:
// - в пам'яті
// - в PostgreSQL
// - у файлі
//
// Він просто працює через інтерфейс BookingStore.
type Service struct {
	store BookingStore
}

// NewService створює новий Service.
//
// store - будь-який об'єкт, який реалізує BookingStore.
func NewService(store BookingStore) *Service {
	return &Service{
		store: store,
	}
}

// Book створює бронювання.
//
// Зараз метод просто передає виклик
// до сховища (store).
func (s *Service) Book(b Booking) error {
	return s.store.Book(b)
}
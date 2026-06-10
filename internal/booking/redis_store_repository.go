package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const defaultHoldTTL = 2 * time.Minute

// інтерфейс RedisStore реалізує систему бронювання місць на основі сесій,
// де дані зберігаються в Redis.

// Схема ключів:

// seat:{movieID}:{seatID}
//     → JSON із даними сесії
//       (TTL встановлено = місце тимчасово утримується,
//        TTL немає = бронювання підтверджене)

// session:{sessionID}
//     → ключ місця
//       (зворотний пошук)
	// RedisStore - реалізація BookingStore,
	// яка зберігає дані в Redis.
type RedisStore struct {
	redisDataBase *redis.Client
}

// NewRedisStore створює новий RedisStore
// та зберігає клієнт Redis для подальшої роботи.
func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{redisDataBase: rdb}
}

// sessionKey формує ключ для пошуку сесії в Redis.
//
// Приклад:
// sessionKey("abc123")
//
// поверне:
//
// session:abc123
func sessionKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}

// Book створює бронювання.
//
// Поки що метод лише викликає hold(),
// який тимчасово резервує місце на певний час (TTL).
func (s *RedisStore) Book(b Booking) error {
	session, err := s.hold(b)
	if err != nil {
		return err
	}

	log.Printf("Сесія заброньована: %+v", session)

	return nil
}

// ListBookings повинен повертати список бронювань.
//
// Наразі ще не реалізований.
func (s *RedisStore) ListBookings(movieID string) []Booking {
	return []Booking{}
}

// hold тимчасово резервує місце.
//
// Алгоритм:
// 1. Генеруємо унікальний sessionID.
// 2. Створюємо Redis-ключ для місця.
// 3. Записуємо бронювання в Redis з TTL.
// 4. Якщо місце вже зайняте - повертаємо помилку.
// 5. Створюємо зворотний запис session -> seat.
// 6. Повертаємо інформацію про створену сесію.
func (s *RedisStore) hold(b Booking) (Booking, error) {

	// Генеруємо унікальний ідентифікатор сесії.
	id := uuid.New().String()

	// Поточний час.
	now := time.Now()

	// Контекст для роботи з Redis.
	ctx := context.Background()

	// Формуємо ключ місця.
	//
	// Наприклад:
	//
	// seat:movie-1:A5
	key := fmt.Sprintf("seat:%s:%s", b.MovieID, b.SeatID)

	// Зберігаємо ID у бронюванні.
	b.ID = id

	// Перетворюємо структуру Booking у JSON,
	// щоб записати її в Redis.
	val, _ := json.Marshal(b)

	// Записуємо значення в Redis.
	//
	// NX означає:
	// "створити ключ лише якщо його ще не існує".
	//
	// Якщо місце вже заброньоване,
	// запис не буде створений.
	response := s.redisDataBase.SetArgs(ctx, key, val, redis.SetArgs{
		Mode: "NX",
		TTL:  defaultHoldTTL,
	})

	// Redis повертає "OK",
	// якщо запис був успішно створений.
	ok := response.Val() == "OK"

	if !ok {
		return Booking{}, ErrSeatsAlreadyBooked
	}

	// Створюємо зворотний зв'язок:
	//
	// session:abc123
	//     ↓
	// seat:movie-1:A5
	//
	// Це дозволить швидко знайти місце
	// за sessionID.
	s.redisDataBase.Set(
		ctx,
		sessionKey(id),
		key,
		defaultHoldTTL,
	)

	// Повертаємо інформацію
	// про створене бронювання.
	return Booking{
		ID:        id,
		MovieID:   b.MovieID,
		SeatID:    b.SeatID,
		UserID:    b.UserID,
		Status:    "held",
		ExpiresAt: now.Add(defaultHoldTTL),
	}, nil
}
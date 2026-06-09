package booking

import "errors"

var (
	ErrSeatsAlreadyBooked = errors.New("Seats are already booked")
)

type Booking struct {
	ID string
	MovieID string
	SeatID string
	UserID string
	Status string
}

type  BookingStore interface {
	Book(b Booking) error
	ListBookings(movieID string) []Booking
}
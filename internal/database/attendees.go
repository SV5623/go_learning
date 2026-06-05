package database

import (
	"context"
	"database/sql"
	"time"
)

type AttendeeModel struct {
	DB *sql.DB
}

type Attendee struct {
	ID int `json:"id"`
	EventID int `json:"event_id"`
	UserID int `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (am *AttendeeModel) AddAttendee(eventId int, userId int) error {
	query := "INSERT INTO attendees (event_id, user_id) VALUES ($1, $2)"
	_, err := am.DB.Exec(query, eventId, userId)
	return err
}	

func (am *AttendeeModel) GetByEventAndAttendee(eventId int, userId int) (*Attendee, error) {
	attendee := &Attendee{}
	query := "SELECT id, event_id, user_id FROM attendees WHERE event_id = $1 AND user_id = $2"
	err := am.DB.QueryRow(query, eventId, userId).Scan(&attendee.ID, &attendee.EventID, &attendee.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return attendee, nil
}

func (am *AttendeeModel) GetEventsByAttendee(attendeeId int) ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	query := `
		SELECT e.id, e.owner_id, e.title, e.description, e.start_time, e.end_time
		FROM events e
		JOIN attendees a ON e.id = a.event_id
		WHERE a.user_id = $1
	`
	rows, err := am.DB.QueryContext(ctx, query, attendeeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.OwnerId, &event.Title, &event.Description, &event.StartTime, &event.EndTime)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (am *AttendeeModel) Get(id int) (*Attendee, error) {
	attendee := &Attendee{}
	query := "SELECT id, event_id, user_id FROM attendees WHERE id = $1"
	err := am.DB.QueryRow(query, id).Scan(&attendee.ID, &attendee.EventID, &attendee.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return attendee, nil
}

func (am *AttendeeModel) Insert(attendee *Attendee) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO attendees (event_id, user_id) VALUES ($1, $2) RETURNING id"
	err := am.DB.QueryRowContext(ctx, query, attendee.EventID, attendee.UserID).Scan(&attendee.ID)
	if err != nil {
		return nil, err
	}
	
	return attendee, nil
}

func (am *AttendeeModel) Delete(eventId int, userId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	query := "DELETE FROM attendees WHERE event_id = $1 AND user_id = $2"
	_, err := am.DB.ExecContext(ctx, query, eventId, userId)
	return err
}

func (am *AttendeeModel) GetAttendeesByEventId(eventId int) ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	query := `
		SELECT u.id, u.name, u.email
		FROM attendees a
		JOIN users u ON a.user_id = u.id
		WHERE a.event_id = $1
	`
	rows, err := am.DB.QueryContext(ctx, query, eventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}


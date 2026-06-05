// Repository: EventRepository
package database

import (
	"context"
	"database/sql"
	"time"
)

type EventModel struct {
	DB *sql.DB
}

type Event struct {
	ID int `json:"id"`
	OwnerId int `json:"owner_id"`
	Title string `json:"title" binding:"min=3,max=100"`
	Description string `json:"description"`
	StartTime string `json:"start_time"`
	EndTime string `json:"end_time"`
	Location string `json:"location"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}	

func (em *EventModel) Create(event *Event) error {
	query := "INSERT INTO events (owner_id, title, description, start_time, end_time, location) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	return em.DB.QueryRow(query, event.OwnerId, event.Title, event.Description, event.StartTime, event.EndTime, event.Location).Scan(&event.ID)
}

func (em *EventModel) Get(id int) (*Event, error) {
	event := &Event{}
	query := "SELECT id, owner_id, title, description, start_time, end_time, location FROM events WHERE id = $1"
	err := em.DB.QueryRow(query, id).Scan(&event.ID, &event.OwnerId, &event.Title, &event.Description, &event.StartTime, &event.EndTime, &event.Location)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (em *EventModel) GetAll() ([]*Event, error) {
	query := "SELECT id, owner_id, title, description, start_time, end_time, location FROM events"
	rows, err := em.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		event := &Event{}
		if err := rows.Scan(&event.ID, &event.OwnerId, &event.Title, &event.Description, &event.StartTime, &event.EndTime, &event.Location); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (em *EventModel) Update(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "UPDATE events SET owner_id = $1, title = $2, description = $3, start_time = $4, end_time = $5, location = $6, updated_at = NOW() WHERE id = $7"
	_, err := em.DB.ExecContext(ctx, query, event.OwnerId, event.Title, event.Description, event.StartTime, event.EndTime, event.Location, event.ID)
	return err
}

func (em *EventModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "DELETE FROM events WHERE id = $1"
	_, err := em.DB.ExecContext(ctx, query, id)
	return err
}


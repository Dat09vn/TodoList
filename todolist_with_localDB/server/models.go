package main

import "time"

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"New"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

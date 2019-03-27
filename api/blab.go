package main

import (
	"encoding/json"
	"time"
)

// User stores info about a user
type User struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Blab stores a blab and its details
type Blab struct {
	ID       uint32    `json:"id"`
	PostTime time.Time `json:"postTime"`
	Author   User      `json:"author"`
	Message  string    `json:"message"`
}

// MarshalJSON converts a Blab into JSON
func (blab Blab) MarshalJSON() ([]byte, error) {
	defaultBlab := struct {
		ID       uint32 `json:"id"`
		PostTime int64  `json:"postTime"`
		Author   User   `json:"author"`
		Message  string `json:"message"`
	}{
		ID:       blab.ID,
		PostTime: blab.PostTime.Unix(),
		Author:   blab.Author,
		Message:  blab.Message,
	}

	return json.Marshal(defaultBlab)
}

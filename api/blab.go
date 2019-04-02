package main

import (
	"encoding/json"
	"strconv"
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

// MarshalJSON converts a User into JSON
func (user User) MarshalJSON() ([]byte, error) {
	defaultBlab := struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		Name:  user.Name,
		Email: user.Email,
	}

	return json.Marshal(defaultBlab)
}

// MarshalJSON converts a Blab into JSON
func (blab Blab) MarshalJSON() ([]byte, error) {
	defaultBlab := struct {
		ID       string `json:"id"`
		PostTime int64  `json:"postTime"`
		Author   User   `json:"author"`
		Message  string `json:"message"`
	}{
		ID:       strconv.FormatUint(uint64(blab.ID), 10),
		PostTime: blab.PostTime.Unix(),
		Author:   blab.Author,
		Message:  blab.Message,
	}

	return json.Marshal(defaultBlab)
}

package main

import (
  "encoding/json"
  "time"
)

// User

type User struct {
  Id uint32 `json:"id"`
  Name string `json:"name"`
  Email string `json:"email"`
}

// Blab

type Blab struct {
  Id uint32 `json:"id"`
  PostTime time.Time `json:"postTime"`
  Author User `json:"author"`
  Message string `json:"message"`
}

func (blab Blab) MarshalJSON() ([]byte, error) {
	defaultBlab := struct {
    Id uint32 `json:"id"`
    PostTime int64 `json:"postTime"`
    Author User `json:"author"`
    Message string `json:"message"`
	}{
		Id: blab.Id,
    PostTime: blab.PostTime.Unix(),
		Author: blab.Author,
		Message: blab.Message,
	}

	return json.Marshal(defaultBlab)
}

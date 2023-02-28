package types

import "time"

type User struct {
	Uid            string    `json:"-"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	LastLogin      time.Time `json:"-"`
	ProfilePicture string    `json:"profilePicture"`
}

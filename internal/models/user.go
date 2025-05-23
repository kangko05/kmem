package models

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u User) String() string {
	return fmt.Sprintf("username: %s\npassword: %s\n", u.Username, u.Password)
}

func (u User) Json() ([]byte, error) {
	return json.Marshal(User{Username: u.Username, Password: u.Password})
}

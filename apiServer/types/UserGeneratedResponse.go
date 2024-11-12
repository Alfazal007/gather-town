package types

import (
	"github.com/Alfazal007/gather-town/internal/database"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func ReturnCreatedUser(userFromDatabase database.User) User {
	return User{
		ID:       userFromDatabase.ID.String(),
		Username: userFromDatabase.Username,
		Email:    userFromDatabase.Email,
	}
}

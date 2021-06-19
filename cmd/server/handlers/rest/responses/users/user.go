package responses

import (
	"github.com/egsam98/voting/gosuslugi/db/queriesdb/usersdb"
)

type User struct {
	ID        int64  `json:"id"`
	Passport  string `json:"passport"`
	Fullname  string `json:"fullname"`
	BirthDate int64  `json:"birth_date"`
	DeathDate *int64 `json:"death_date"`
}

func NewUser(user usersdb.User) User {
	res := User{
		ID:        user.ID,
		Passport:  user.Passport,
		Fullname:  user.Fullname,
		BirthDate: user.BirthDate.Unix(),
	}
	if user.DeathDate.Valid {
		res.DeathDate = new(int64)
		*res.DeathDate = user.DeathDate.Time.Unix()
	}
	return res
}

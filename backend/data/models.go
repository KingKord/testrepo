package data

import "database/sql"

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: User{},
	}
}

type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic,omitempty"`
	Agify       string `json:"agify,omitempty"`
	Genderize   string `json:"genderize,omitempty"`
	Nationalize string `json:"nationalize"`
}

type Models struct {
	User User
}

func GetUserByName(name string) User {

	return User{}
}

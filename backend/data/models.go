package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var db *sql.DB

const dbTimeout = time.Second * 3

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
	Agify       int    `json:"agify,omitempty"`
	Genderize   string `json:"genderize,omitempty"`
	Nationalize string `json:"nationalize"`
}

type Models struct {
	User User
}

func (u *User) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `insert into users (name, surname, patronymic, agify, genderize, nationalize)
		values ($1, $2, $3, $4, $5, $6) returning id`

	err := db.QueryRowContext(ctx, stmt,
		user.Name,
		user.Surname,
		user.Patronymic,
		string(rune(user.Agify)),
		user.Genderize,
		user.Nationalize,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}
	u.ID = newID

	return newID, nil
}

func GetInfoFromOpenAPI(URL string) (*http.Response, error) {
	request, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("can't get info from API")
	}

	fmt.Println("URL:", URL)
	fmt.Println("Response Status:", response.Status)

	return response, nil
}

func GetUserByName(name string) User {

	return User{}
}

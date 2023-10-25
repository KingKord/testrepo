package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
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
	Age         int    `json:"age,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Nationality string `json:"nationality,omitempty"`
}

type Models struct {
	User User
}

func (u *User) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `insert into users (name, surname, patronymic, age, gender, nationality)
		values ($1, $2, $3, $4, $5, $6) returning id`

	err := db.QueryRowContext(ctx, stmt,
		user.Name,
		user.Surname,
		user.Patronymic,
		user.Age,
		user.Gender,
		user.Nationality,
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

	return response, nil
}

func (u *User) GetAll(perPage, page int) ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, name, surname, patronymic, age, gender, nationality
	from users order by id limit $1 offset $2`

	offset := (page - 1) * perPage

	rows, err := db.QueryContext(ctx, query, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Surname,
			&user.Patronymic,
			&user.Age,
			&user.Gender,
			&user.Nationality,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (u *User) GetAllUsersByGender(gender string, perPage, page int) ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, name, surname, patronymic, age, gender, nationality
	from users where gender = $1 order by id limit $2 offset $3`
	offset := (page - 1) * perPage

	rows, err := db.QueryContext(ctx, query, gender, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Surname,
			&user.Patronymic,
			&user.Age,
			&user.Gender,
			&user.Nationality,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (u *User) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil

}

func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	//id, name, surname, patronymic, age, gender, nationality
	stmt := `update users set
		name = $1,
		surname = $2,
		patronymic = $3,
		age = $4,
		gender = $5,
		nationality = $6
		where id = $7
	`

	_, err := db.ExecContext(ctx, stmt,
		u.Name,
		u.Surname,
		u.Patronymic,
		u.Age,
		u.Gender,
		u.Nationality,
		u.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

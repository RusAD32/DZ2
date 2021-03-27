package dealership

import (
	"database/sql"
	"fmt"
	"log"

	"DZ2/internal/app/models"
)

type UserRepository struct {
	store *Store
}

var (
	tableUser string = "users"
)

//Create user in database
func (ur *UserRepository) Create(u *models.User) (*models.User, error) {
	query := fmt.Sprintf("INSERT INTO %s (login, password) VALUES ($1, $2) RETURNING login", tableUser)
	rows, err := ur.store.db.Query(
		query,
		u.Login,
		u.Password,
	)
	if err != nil {
		log.Println(u, err)
		return nil, err
	}
	defer rows.Close()
	return u, nil
}

//Select All
func (ur *UserRepository) FindByLogin(login string) (*models.User, bool, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE login=$1", tableUser)
	row := ur.store.db.QueryRow(query, login)
	u := models.User{}
	err := row.Scan(&u.Login, &u.Password)
	log.Println(u)
	if err == sql.ErrNoRows {
		log.Println(err)
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	return &u, true, nil
}

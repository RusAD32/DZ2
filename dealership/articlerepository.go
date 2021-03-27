package dealership

import (
	"database/sql"
	"fmt"
	"log"

	"DZ2/internal/app/models"
)

type ArticleRepository struct {
	store *Store
}

var (
	carTebleName string = "cars"
)

//For Post request
func (cr *ArticleRepository) Create(a *models.Car) (*models.Car, bool, error) {
	_, found, err := cr.FindCarById(a.Mark)
	if err != nil {
		return nil, false, err
	}
	if found { // don't create cars if we already have one
		return nil, false, nil
	}
	query := fmt.Sprintf("INSERT INTO %s (mark, max_speed, distance, handler, stock) VALUES ($1, $2, $3, $4, $5)", carTebleName)
	rows, err := cr.store.db.Query(query, a.Mark, a.MaxSpeed, a.Distance, a.Handler, a.Stock)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()
	return a, true, nil
}

//For Put request
func (cr *ArticleRepository) Update(a *models.Car) (*models.Car, bool, error) {
	_, found, err := cr.FindCarById(a.Mark)
	if err != nil {
		return nil, false, err
	}
	if !found { // don't update cars we don't have
		return nil, false, nil
	}
	query := fmt.Sprintf("UPDATE %s SET max_speed=$2, distance=$3, handler=$4, stock=$5 WHERE mark=$1", carTebleName)
	rows, err := cr.store.db.Query(query, a.Mark, a.MaxSpeed, a.Distance, a.Handler, a.Stock)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()
	return a, true, nil
}

//For DELETE request
func (cr *ArticleRepository) DeleteById(mark string) (*models.Car, bool, error) {
	article, ok, err := cr.FindCarById(mark)
	if err != nil {
		return nil, false, err
	}
	if ok {
		query := fmt.Sprintf("delete from %s where mark=$1", carTebleName)
		_, err = cr.store.db.Exec(query, mark)
		if err != nil {
			return nil, false, err
		}
	}

	return article, ok, nil
}

//Helper for Delete by id and GET by id request
func (cr *ArticleRepository) FindCarById(mark string) (*models.Car, bool, error) {
	car, err := cr.SelectOne(mark)
	if err == sql.ErrNoRows {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	return car, true, nil
}

//Get all request and helper for FindByID
func (cr *ArticleRepository) SelectAll() ([]*models.Car, error) {
	query := fmt.Sprintf("SELECT * FROM %s", carTebleName)
	rows, err := cr.store.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cars := make([]*models.Car, 0)
	for rows.Next() {
		a := models.Car{}
		err := rows.Scan(&a.Mark, &a.MaxSpeed, &a.Distance, &a.Handler, &a.Stock)
		if err != nil {
			log.Println(err)
			continue
		}
		cars = append(cars, &a)
	}
	return cars, nil
}

//Get all request and helper for FindByID
func (cr *ArticleRepository) SelectOne(id string) (*models.Car, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE mark=$1", carTebleName)
	row := cr.store.db.QueryRow(query, id)
	a := models.Car{}
	err := row.Scan(&a.Mark, &a.MaxSpeed, &a.Distance, &a.Handler, &a.Stock)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

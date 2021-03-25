package dealership

import (
	"fmt"
	"log"

	"DZ2/internal/app/models"
)

type ArticleRepository struct {
	store *Store
}

var (
	tableArticle string = "cars"
)

//For Post request
func (ar *ArticleRepository) Create(a *models.Car) (*models.Car, bool, error) {
	_, found, err := ar.FindArticleById(a.Mark)
	if err != nil {
		return nil, false, err
	}
	if found { // don't create cars if we already have one
		return nil, false, nil
	}
	query := fmt.Sprintf("INSERT INTO %s (mark, max_speed, distance, handler, stock) VALUES ($1, $2, $3, $4, $5)", tableArticle)
	if err := ar.store.db.QueryRow(query, a.Mark, a.MaxSpeed, a.Distance, a.Handler, a.Stock).Scan(); err != nil {
		return nil, false, err
	}
	return a, true, nil
}

//For DELETE request
func (ar *ArticleRepository) DeleteById(mark string) (*models.Car, error) {
	article, ok, err := ar.FindArticleById(mark)
	if err != nil {
		return nil, err
	}
	if ok {
		query := fmt.Sprintf("delete from %s where mark=$1", tableArticle)
		_, err = ar.store.db.Exec(query, mark)
		if err != nil {
			return nil, err
		}
	}

	return article, nil
}

//Helper for Delete by id and GET by id request
func (ar *ArticleRepository) FindArticleById(id string) (*models.Car, bool, error) {
	articles, err := ar.SelectAll()
	founded := false
	if err != nil {
		return nil, founded, err
	}
	var articleFinded *models.Car
	for _, a := range articles {
		if a.ID == id {
			articleFinded = a
			founded = true
		}
	}

	return articleFinded, founded, nil

}

//Get all request and helper for FindByID
func (ar *ArticleRepository) SelectAll() ([]*models.Car, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableArticle)
	rows, err := ar.store.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	articles := make([]*models.Car, 0)
	for rows.Next() {
		a := models.Car{}
		err := rows.Scan(&a.ID, &a.Title, &a.Author, &a.Content)
		if err != nil {
			log.Println(err)
			continue
		}
		articles = append(articles, &a)
	}
	return articles, nil
}

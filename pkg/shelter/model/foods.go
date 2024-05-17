package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Food struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Quantity string `json:"quantity"`
	Shelter  string `json:"shelter"`
}

type FoodModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (e FoodModel) Insert(employee *Food) error {
	// check for ID needed here if error
	query := `
		INSERT INTO Foods (Name, Type, Quantity,Shelter) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, Name
		`
	// check if its animal of Animals in case of error
	args := []interface{}{employee.Name, employee.Type, employee.Quantity, employee.Shelter}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return e.DB.QueryRowContext(ctx, query, args...).Scan(&employee.ID, &employee.Name)
}

func (f FoodModel) Get(id int) (*Food, error) {
	// Retrieve a specific menu item based on its ID.
	query := `
		SELECT id, Name, Type, Quantity, Shelter
		FROM Foods
		WHERE ID = $1
		`
	var food Food
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// again animal or Animals?
	row := f.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&food.ID, &food.Name, &food.Type, &food.Quantity, &food.Shelter)
	if err != nil {
		return nil, err
	}
	return &food, nil
}

func (f FoodModel) GetSort(Name, Type string, filters Filters) ([]*Food, Metadata, error) {

	// Retrieve all menu items from the database.
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, Name, Type, Quantity, Shelter
		FROM Foods
		WHERE (LOWER(Name) = LOWER($1) OR $1 = '')
		AND (LOWER(Type) = LOWER($2) OR $2 = '')
		--AND (Quantity >= $2 OR $2 = 0)
		--AND (Quantity <= $3 OR $3 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4		`,
		filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Organize our four placeholder parameter values in a slice.
	args := []interface{}{Name, Type, filters.limit(), filters.offset()}

	// log.Println(query, title, from, to, filters.limit(), filters.offset())
	// Use QueryContext to execute the query. This returns a sql.Rows result set containing
	// the result.
	rows, err := f.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	// Importantly, defer a call to rows.Close() to ensure that the result set is closed
	// before GetAll returns.
	defer func() {
		if err := rows.Close(); err != nil {
			f.ErrorLog.Println(err)
		}
	}()

	// Declare a totalRecords variable
	totalRecords := 0

	var foods []*Food
	for rows.Next() {
		var food Food
		err := rows.Scan(&totalRecords, &food.ID, &food.Name, &food.Type, &food.Quantity, &food.Shelter)
		if err != nil {
			return nil, Metadata{}, err
		}

		// Add the Movie struct to the slice
		foods = append(foods, &food)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Generate a Metadata struct, passing in the total record count and pagination parameters
	// from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// If everything went OK, then return the slice of the movies and metadata.
	return foods, metadata, nil
}

func (f FoodModel) Delete(id int) error {
	// Delete a specific menu item from the database.
	query := `
		DELETE FROM Foods
		WHERE ID = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := f.DB.ExecContext(ctx, query, id)
	return err
}

func (f FoodModel) Update(food *Food) error {
	// Update a specific animal in the database.
	query := `
        UPDATE Foods
        SET Name = $2, Type = $3, Quantity = $4, Shelter = $5
        WHERE ID = $1
        RETURNING ID
        `
	args := []interface{}{food.ID, food.Name, food.Type, food.Quantity, food.Shelter}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return f.DB.QueryRowContext(ctx, query, args...).Scan(&food.ID)
}

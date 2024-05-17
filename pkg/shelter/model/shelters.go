package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Shelter struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
	Capacity    string `json:"capacity"`
}

type ShelterModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (s ShelterModel) Insert(shelter *Shelter) error {
	query := `
    INSERT INTO Shelters (Name, Location, Description, Capacity) 
    VALUES ($1, $2, $3, $4) 
    RETURNING ID, Name
  `
	args := []interface{}{shelter.Name, shelter.Location, shelter.Description, shelter.Capacity}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&shelter.ID, &shelter.Name)
}

func (s ShelterModel) Get(id int) (*Shelter, error) {
	query := `
    SELECT ID, Name, Location, Description, Capacity
    FROM Shelters
    WHERE ID = $1
  `
	var shelter Shelter
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := s.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&shelter.ID, &shelter.Name, &shelter.Location, &shelter.Description, &shelter.Capacity)
	if err != nil {
		return nil, err
	}
	return &shelter, nil
}

func (s ShelterModel) GetSort(location, name string, filters Filters) ([]*Shelter, Metadata, error) {
	query := fmt.Sprintf(`
    SELECT count(*) OVER(), ID, Name, Location, Description, Capacity
    FROM Shelters
    WHERE (LOWER(Location) = LOWER($1) OR $1 = '')
    AND (LOWER(Name) = LOWER($2) OR $2 = '')
    ORDER BY %s %s, ID ASC
    LIMIT $3 OFFSET $4
  `, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{location, name, filters.limit(), filters.offset()}

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			s.ErrorLog.Println(err)
		}
	}()

	totalRecords := 0
	var shelters []*Shelter
	for rows.Next() {
		var shelter Shelter
		err := rows.Scan(&totalRecords, &shelter.ID, &shelter.Name, &shelter.Location, &shelter.Description, &shelter.Capacity)
		if err != nil {
			return nil, Metadata{}, err
		}
		shelters = append(shelters, &shelter)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return shelters, metadata, nil
}

func (s ShelterModel) Delete(id int) error {
	query := `
  DELETE FROM Shelters
	WHERE ID = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := s.DB.ExecContext(ctx, query, id)
	return err
}

func (s ShelterModel) Update(shelter *Shelter) error {
	query := `
        UPDATE Shelters
        SET Name = $2, Location = $3, Description = $4, Capacity = $5
        WHERE ID = $1
        RETURNING ID
    `
	args := []interface{}{shelter.ID, shelter.Name, shelter.Location, shelter.Description, shelter.Capacity}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&shelter.ID)
}

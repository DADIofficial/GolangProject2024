package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Volunteer struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Age          string `json:"age"`
	Description  string `json:"description"`
	Role         string `json:"role"`
	Organization string `json:"organization"`
}

type VolunteerModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (v VolunteerModel) Insert(volunteer *Volunteer) error {
	query := `
	INSERT INTO Volunteers (Name, Surname, Age, Description, Role, Organization) 
	  VALUES ($1, $2, $3, $4, $5, $6) 
	  RETURNING ID, Name
	  `

	args := []interface{}{volunteer.Name, volunteer.Surname, volunteer.Age, volunteer.Description, volunteer.Role, volunteer.Organization}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return v.DB.QueryRowContext(ctx, query, args...).Scan(&volunteer.ID, &volunteer.Name)
}

func (v VolunteerModel) Get(id int) (*Volunteer, error) {
	query := `
	SELECT id, Name, Surname, Age, Description, Role, Organization
	FROM Volunteers
	WHERE ID = $1
	  `
	var volunteer Volunteer
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := v.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&volunteer.ID, &volunteer.Name, &volunteer.Surname, &volunteer.Age, &volunteer.Description, &volunteer.Role, &volunteer.Organization)
	if err != nil {
		return nil, err
	}
	return &volunteer, nil
}

func (v VolunteerModel) GetSort(Name, Role string, filters Filters) ([]*Volunteer, Metadata, error) {

	query := fmt.Sprintf(
		`
	  SELECT count(*) OVER(), ID, Name, Surname, Age, Description, Role, Organization
	  FROM Volunteers
	  WHERE (LOWER(Name) = LOWER($1) OR $1 = '')
	  AND (LOWER(Role) = LOWER($2) OR $2 = '')
	  ORDER BY %s %s, id ASC
	  LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{Name, Role, filters.limit(), filters.offset()}

	rows, err := v.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			v.ErrorLog.Println(err)
		}
	}()

	totalRecords := 0

	var volunteers []*Volunteer
	for rows.Next() {
		var volunteer Volunteer
		err := rows.Scan(&totalRecords, &volunteer.ID, &volunteer.Name, &volunteer.Surname, &volunteer.Age, &volunteer.Description, &volunteer.Role, &volunteer.Organization)
		if err != nil {
			return nil, Metadata{}, err
		}
		volunteers = append(volunteers, &volunteer)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return volunteers, metadata, nil
}

func (v VolunteerModel) Delete(id int) error {
	query := `
	  DELETE FROM Volunteers
	  WHERE ID = $1
	  `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := v.DB.ExecContext(ctx, query, id)
	return err
}

func (v VolunteerModel) Update(volunteer *Volunteer) error {
	query := `
	  UPDATE Volunteers
	  SET Name = $2, Surname = $3, Age = $4, Description = $5, Role = $6, Organization = $7
	  WHERE ID = $1
	  RETURNING ID
		  `
	args := []interface{}{volunteer.ID, volunteer.Name, volunteer.Surname, volunteer.Age, volunteer.Description, volunteer.Role, volunteer.Organization}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return v.DB.QueryRowContext(ctx, query, args...).Scan(&volunteer.ID)
}

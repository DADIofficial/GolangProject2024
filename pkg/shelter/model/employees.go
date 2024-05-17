package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Employee struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Salary  string `json:"salary"`
	Duty    string `json:"duty"`
	Shelter string `json:"shelter"`
}

type EmployeeModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (e EmployeeModel) Insert(employee *Employee) error {
	// check for ID needed here if error
	query := `
		INSERT INTO Employees (Name, Surname, Salary, Duty, Shelter) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING ID, Name
		`
	// check if its animal of Animals in case of error
	args := []interface{}{employee.Name, employee.Surname, employee.Salary, employee.Duty, employee.Shelter}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return e.DB.QueryRowContext(ctx, query, args...).Scan(&employee.ID, &employee.Name)
}

func (e EmployeeModel) Get(id int) (*Employee, error) {
	// Retrieve a specific menu item based on its ID.
	query := `
		SELECT ID, Name, Surname, Salary, Duty, Shelter
		FROM Employees
		WHERE ID = $1
		`
	var employee Employee
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// again animal or Animals?
	row := e.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&employee.ID, &employee.Name, &employee.Surname, &employee.Salary, &employee.Duty, &employee.Shelter)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (e EmployeeModel) GetSort(Name, Shelter string, filters Filters) ([]*Employee, Metadata, error) {

	// Retrieve all menu items from the database.
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, Name, Surname, Salary, Duty, Shelter
		FROM Employees
		WHERE (LOWER(Name) = LOWER($1) OR $1 = '')
		AND (LOWER(Shelter) = LOWER($2) OR $2 = '')
		--AND (Salary >= $2 OR $2 = 0)
		--AND (Salary <= $3 OR $3 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4		`,
		filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Organize our four placeholder parameter values in a slice.
	args := []interface{}{Name, Shelter, filters.limit(), filters.offset()}

	// log.Println(query, title, from, to, filters.limit(), filters.offset())
	// Use QueryContext to execute the query. This returns a sql.Rows result set containing
	// the result.
	rows, err := e.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	// Importantly, defer a call to rows.Close() to ensure that the result set is closed
	// before GetAll returns.
	defer func() {
		if err := rows.Close(); err != nil {
			e.ErrorLog.Println(err)
		}
	}()

	// Declare a totalRecords variable
	totalRecords := 0

	var employees []*Employee
	for rows.Next() {
		var employee Employee
		err := rows.Scan(&totalRecords, &employee.ID, &employee.Name, &employee.Surname, &employee.Salary, &employee.Duty, &employee.Shelter)
		if err != nil {
			return nil, Metadata{}, err
		}

		// Add the Movie struct to the slice
		employees = append(employees, &employee)
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
	return employees, metadata, nil
}

func (e EmployeeModel) Delete(id int) error {
	// Delete a specific menu item from the database.
	query := `
		DELETE FROM Employees
		WHERE ID = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := e.DB.ExecContext(ctx, query, id)
	return err
}

func (e EmployeeModel) Update(employee *Employee) error {
	// Update a specific animal in the database.
	query := `
        UPDATE Employees
        SET Name = $2, Surname = $3, Salary = $4, Duty = $5, Shelter = $6
        WHERE ID = $1
        RETURNING ID
        `
	args := []interface{}{employee.ID, employee.Name, employee.Surname, employee.Salary, employee.Duty, employee.Shelter}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return e.DB.QueryRowContext(ctx, query, args...).Scan(&employee.ID)
}

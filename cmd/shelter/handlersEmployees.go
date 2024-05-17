package main

import (
	//"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/po133na/go-mid/pkg/shelter/model"
	"github.com/po133na/go-mid/pkg/shelter/validator"
)

func (app *application) createEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Surname string `json:"surname"`
		Salary  string `json:"salary"`
		Duty    string `json:"duty"`
		Shelter string `json:"shelter"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	employee := &model.Employee{
		ID:      input.ID,
		Name:    input.Name,
		Surname: input.Surname,
		Salary:  input.Salary,
		Duty:    input.Duty,
		Shelter: input.Shelter,
	}

	err = app.models.Employees.Insert(employee)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, employee)
}

func (app *application) getEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["employeeId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Employee ID")
		return
	}

	employee, err := app.models.Employees.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, employee)
}

func (app *application) getEmployeesSortedHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name    string `json:"name"`
		Surname string `json:"surname"`
		model.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Name = app.readStrings(qs, "name", "")
	input.Surname = app.readStrings(qs, "surname", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readStrings(qs, "sort", "id")

	input.Filters.SortSafeList = []string{
		// ascending sort values
		"id", "name", "surname", "salary", "duty", "shelter",
		// descending sort values
		"-id", "-name", "-surname", "-salary", "-duty", "-shelter",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	employees, metadata, err := app.models.Employees.GetSort(input.Name, input.Surname, input.Filters)
	if err != nil {
		fmt.Println("We are in search employees handler", "\nname: ", input.Name, "\nsurname:", input.Surname, "\n", input.Filters)
		fmt.Print("\nError: ", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"employees": employees, "metadata": metadata}, nil)
}

func (app *application) updateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["employeeId"] // CHECK HERE FOR ERRORS

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Employee ID")
		return
	}

	employee, err := app.models.Employees.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Name   *string `json:"name"`
		Salary *string `json:"salary"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Name != nil {
		employee.Name = *input.Name
	}

	if input.Salary != nil {
		employee.Salary = *input.Salary
	}
	err = app.models.Employees.Update(employee)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) deleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["employeeId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Employee ID")
		return
	}

	err = app.models.Employees.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

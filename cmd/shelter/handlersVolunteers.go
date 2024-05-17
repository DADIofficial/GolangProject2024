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

func (app *application) createVolunteerHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Surname      string `json:"surname"`
		Age          string `json:"age"`
		Description  string `json:"description"`
		Role         string `json:"role"`
		Organization string `json:"organization"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	volunteer := &model.Volunteer{
		ID:           input.ID,
		Name:         input.Name,
		Surname:      input.Surname,
		Age:          input.Age,
		Description:  input.Description,
		Role:         input.Role,
		Organization: input.Organization,
	}

	err = app.models.Volunteers.Insert(volunteer)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, volunteer)
}

func (app *application) getVolunteerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["volunteerId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Volunteer ID")
		return
	}

	volunteer, err := app.models.Volunteers.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, volunteer)
}

func (app *application) getVolunteersSortedHandler(w http.ResponseWriter, r *http.Request) {

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
		"id", "name", "surname", "age", "description", "role", "organization",
		// descending sort values
		"-id", "-name", "-surname", "-age", "-description", "-role", "-organization",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	volunteers, metadata, err := app.models.Volunteers.GetSort(input.Name, input.Surname, input.Filters)
	if err != nil {
		fmt.Println("We are in search volunteers handler", "\nname: ", input.Name, "\nsurname:", input.Surname, "\n", input.Filters)
		fmt.Print("\nError: ", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"volunteers": volunteers, "metadata": metadata}, nil)
}

func (app *application) updateVolunteerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["volunteerId"] // CHECK HERE FOR ERRORS

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Volunteer ID")
		return
	}

	volunteer, err := app.models.Volunteers.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Name != nil {
		volunteer.Name = *input.Name
	}

	if input.Description != nil {
		volunteer.Description = *input.Description
	}
	err = app.models.Volunteers.Update(volunteer)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) deleteVolunteerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["volunteerId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Volunteer ID")
		return
	}

	err = app.models.Volunteers.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

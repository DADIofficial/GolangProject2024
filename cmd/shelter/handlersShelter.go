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

func (app *application) createShelterHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Location    string `json:"location"`
		Description string `json:"description"`
		Capacity    string `json:"capacity"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	shelter := &model.Shelter{
		ID:          input.ID,
		Name:        input.Name,
		Location:    input.Location,
		Description: input.Description,
		Capacity:    input.Capacity,
	}

	err = app.models.Shelters.Insert(shelter)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, shelter)
}

func (app *application) getShelterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["shelterId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid shelter ID")
		return
	}

	shelter, err := app.models.Shelters.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, shelter)
}

func (app *application) getSheltersSortedHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// param := vars["animalId"]
	var input struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		model.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Name = app.readStrings(qs, "name", "")
	input.Location = app.readStrings(qs, "location", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readStrings(qs, "sort", "id")

	input.Filters.SortSafeList = []string{
		// ascending sort values
		"id", "name", "location", "description", "capacity",
		// descending sort values
		"-id", "-name", "-location", "-description", "-capacity",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	shelters, metadata, err := app.models.Shelters.GetSort(input.Name, input.Location, input.Filters)
	if err != nil {
		fmt.Println("We are in search shelters handler", "\nshelter: ", input.Name, "\nlocation:", input.Location, "\n", input.Filters)
		fmt.Print("\nError: ", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"shelters": shelters, "metadata": metadata}, nil)
}

func (app *application) updateShelterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["shelterId"] // CHECK HERE FOR ERRORS

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid shelter ID")
		return
	}

	shelter, err := app.models.Shelters.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Name     *string `json:"name"`
		Capacity *string `json:"capacity"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Name != nil {
		shelter.Name = *input.Name
	}

	if input.Capacity != nil {
		shelter.Capacity = *input.Capacity
	}
	err = app.models.Shelters.Update(shelter)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) deleteShelterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["shelterId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Shelter ID")
		return
	}

	err = app.models.Shelters.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

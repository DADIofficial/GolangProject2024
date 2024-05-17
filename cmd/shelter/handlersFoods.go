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

func (app *application) createFoodHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Type     string `json:"type"`
		Quantity string `json:"quantity"`
		Shelter  string `json:"shelter"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	food := &model.Food{
		ID:       input.ID,
		Name:     input.Name,
		Type:     input.Type,
		Quantity: input.Quantity,
		Shelter:  input.Shelter,
	}

	err = app.models.Foods.Insert(food)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, food)
}

func (app *application) getFoodHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["foodId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Food ID")
		return
	}

	food, err := app.models.Foods.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, food)
}

func (app *application) getFoodsSortedHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name string `json:"name"`
		Type string `json:"type"`
		model.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Name = app.readStrings(qs, "name", "")
	input.Type = app.readStrings(qs, "type", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readStrings(qs, "sort", "id")

	input.Filters.SortSafeList = []string{
		// ascending sort values
		"id", "name", "type", "quantity", "shelter",
		// descending sort values
		"-id", "-name", "-type", "-quantity", "shelter",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	foods, metadata, err := app.models.Foods.GetSort(input.Name, input.Type, input.Filters)
	if err != nil {
		fmt.Println("We are in search foods handler", "\nname: ", input.Name, "\ntype:", input.Type, "\n", input.Filters)
		fmt.Print("\nError: ", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"foods": foods, "metadata": metadata}, nil)
}

func (app *application) updateFoodHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["foodId"] // CHECK HERE FOR ERRORS

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Food ID")
		return
	}

	food, err := app.models.Foods.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Name     *string `json:"name"`
		Quantity *string `json:"quantity"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Name != nil {
		food.Name = *input.Name
	}

	if input.Quantity != nil {
		food.Quantity = *input.Quantity
	}
	err = app.models.Foods.Update(food)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) deleteFoodHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["foodId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Food ID")
		return
	}

	err = app.models.Foods.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

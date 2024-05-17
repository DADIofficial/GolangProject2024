package main

import (
	"net/http"

	//new
	"github.com/gorilla/mux"
	//new
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	// Convert app.methodNotAllowedResponse helper to a http.Handler and set it as the custom
	// error handler for 405 Method Not Allowed responses
	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	v1 := r.PathPrefix("/api/v1").Subrouter()
	// Animal Singleton
	v1.HandleFunc("/animals", app.createAnimalHandler).Methods("POST")
	v1.HandleFunc("/animals/{animalId:[0-9]+}", app.getAnimalHandler).Methods("GET")
	v1.HandleFunc("/animals", app.getAnimalsSortedHandler).Methods("GET")
	v1.HandleFunc("/animals/{animalId:[0-9]+}", app.updateAnimalHandler).Methods("PUT")
	v1.HandleFunc("/animals/{animalId:[0-9]+}", app.requirePermissions("animals:read", app.deleteAnimalHandler)).Methods("DELETE")

	v1.HandleFunc("/volunteers", app.createVolunteerHandler).Methods("POST")
	v1.HandleFunc("/volunteers/{volunteerId:[0-9]+}", app.getVolunteerHandler).Methods("GET")
	v1.HandleFunc("/volunteers", app.getVolunteersSortedHandler).Methods("GET")
	v1.HandleFunc("/volunteers/{volunteerId:[0-9]+}", app.updateVolunteerHandler).Methods("PUT")
	v1.HandleFunc("/volunteers/{volunteerId:[0-9]+}", app.requirePermissions("animals:read", app.deleteVolunteerHandler)).Methods("DELETE")

	v1.HandleFunc("/shelters", app.createShelterHandler).Methods("POST")
	v1.HandleFunc("/shelters/{shelterId:[0-9]+}", app.getShelterHandler).Methods("GET")
	v1.HandleFunc("/shelters", app.getSheltersSortedHandler).Methods("GET")
	v1.HandleFunc("/shelters/{shelterId:[0-9]+}", app.updateShelterHandler).Methods("PUT")
	v1.HandleFunc("/shelters/{shelterId:[0-9]+}", app.requirePermissions("animals:megauser", app.deleteShelterHandler)).Methods("DELETE")

	v1.HandleFunc("/employees", app.createEmployeeHandler).Methods("POST")
	v1.HandleFunc("/employees/{employeeId:[0-9]+}", app.getEmployeeHandler).Methods("GET")
	v1.HandleFunc("/employees", app.getEmployeesSortedHandler).Methods("GET")
	v1.HandleFunc("/employees/{employeeId:[0-9]+}", app.updateEmployeeHandler).Methods("PUT")
	v1.HandleFunc("/employees/{employeeId:[0-9]+}", app.requirePermissions("animals:read", app.deleteEmployeeHandler)).Methods("DELETE")

	v1.HandleFunc("/foods", app.createFoodHandler).Methods("POST")
	v1.HandleFunc("/foods/{foodId:[0-9]+}", app.getFoodHandler).Methods("GET")
	v1.HandleFunc("/foods", app.getFoodsSortedHandler).Methods("GET")
	v1.HandleFunc("/foods/{foodId:[0-9]+}", app.updateFoodHandler).Methods("PUT")
	v1.HandleFunc("/foods/{foodId:[0-9]+}", app.requirePermissions("animals:read", app.deleteFoodHandler)).Methods("DELETE")

	users1 := r.PathPrefix("/api/v1").Subrouter()

	users1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	users1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	users1.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")

	return app.authenticate(r)
}

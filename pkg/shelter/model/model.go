package model

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

type Models struct {
	User        UserModel
	Animals     AnimalModel
	Shelters    ShelterModel
	Tokens      TokenModel
	Permissions PermissionModel
	Volunteers  VolunteerModel
	Employees   EmployeeModel
	Foods       FoodModel
}

func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		User: UserModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Animals: AnimalModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Shelters: ShelterModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Employees: EmployeeModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Volunteers: VolunteerModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Foods: FoodModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Tokens: TokenModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Permissions: PermissionModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
	}
}

var (
	ErrRecordNotFound = errors.New("record not found")

	ErrEditConflict = errors.New("edit conflict")
)

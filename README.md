# Animal Shelter Project

## Introduction
The Animal Shelter Project is a comprehensive backend solution developed by our team to assist animal shelters in managing their operations efficiently. This project includes a full-fledged API that enables users to interact with the database to find pets for adoption, access detailed information about the pets, and also provides a platform for pet owners to relinquish their pets if necessary.

## Technologies
- **Programming Language:** Go
- **Database:** PostgreSQL

## Features
- **Pet Adoption:** Users can browse through a detailed list of pets available for adoption.
- **Pet Surrender:** Provides a facility for pet owners to give up their pets to the shelter.
- **Database Management:** Manage detailed information about shelters, users, volunteers, employees, and food supplies.

## Team Members
- Polina Stelmakh (22B030588)
- Anel Tulepbergen (22B030602)
- Alina Amreyeva (22B031240)
- Dossym Ibray (22B030545)
## Database Schema

### Entities and Their Attributes

#### Users
- `UserID` (AUTO_INCREMENT PRIMARY KEY)
- `User_Email` (TEXT)
- `Username` (VARCHAR(30))
- `Password` (Encrypted)
- `Number_of_phone_user` (VARCHAR(50))
- `Role` (Default: "User")

#### Animal
- `ID` (AUTO_INCREMENT PRIMARY KEY)
- `Kind_Of_Animal` (VARCHAR(255))
- `Breed_Of_Animal` (VARCHAR(255))
- `Name` (VARCHAR(255))
- `Age` (INTEGER)
- `Description` (TEXT)

#### Admins
- `AdminID` (AUTO_INCREMENT PRIMARY KEY)
- `Admin_Email` (TEXT)
- `Adminame` (VARCHAR(30))
- `Password` (Encrypted)
- `Number_of_phone_Admin` (VARCHAR(50))
- `Role` (VARCHAR(255) Type: Back_end, Front_end, G_admin, etc.)

#### Shelter
- `AdminID` (AUTO_INCREMENT PRIMARY KEY)
- `CreatedAT` (VARCHAR(30))
- `UpdatedAT` (VARCHAR(30))
- `TITLE` (TEXT)
- `DESCRIPTION` (VARCHAR(50))
- `ADDRESS` (VARCHAR(30))
- `COORDINATES` (TEXT)

#### Food
- `ID` (AUTO_INCREMENT PRIMARY KEY)
- `Name` (VARCHAR(30))
- `Type` (VARCHAR(30))
- `Quantity` (TEXT)
- `Shelter` (VARCHAR(50))


#### Role
- `Permissions` (TEXT)

## Conclusion

This project aims to simplify the operations of animal shelters, making it easier to manage pet adoption and care through a robust backend system.

## Build Instructions
To build the project, navigate to the project directory and run the following command in the integrated terminal:
```bash
go run main.go

POST /animals        - Add a new animal to the shelter.
GET /animals/:id     - Retrieve details of a specific animal.
PUT /animals/:id     - Update details of a specific animal.
DELETE /animals/:id  - Remove an animal from the shelter.


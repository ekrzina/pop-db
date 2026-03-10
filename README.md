# PopDB

## Introduction

PopDB is a lightweight, personal-use database project designed to manage and display population census data through a dynamic web application. The project integrates a relational SQL database with a streamlined frontend and backend architecture.

## Technology Stack

- `SQLLite` - embedded relational database for data storage
- `Python` - lightweight frontend interface for client-side interactions
- `Go` - backend server-side logic using the Gin web framework

See details on each part of the project below.

## Implementation

### SQL Database

The core of the PopDB application is a relational database that stores user population data with the following fields:
- **ID** (unique identifier, mandatory) 
- **Name** (mandatory) 
- **Surname** (mandatory) 
- **Occupation** (optional) 
- **Date of Birth** (mandatory) 
- **Nationality** (mandatory) 
- **City** (mandatory) 
- **Notes** (mandatory) 
- **Picture** (optional)
- **Height** (mandatory)
- **Weight** (mandatory)
- **Blood Type** (mandatory)
- **Medical Conditions** (optional)

The database is stored in a single `.sqlite` file located in the `assets/` directory. On running the application, if the database file is not found, the backend service will automatically initialize and generate the necessary database structure.

### OpenAPI

The OpenAPI is created using `Swagger Open API`, generated with `openapi` using:

```bash
go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest --config api/server/server.cfg.yaml api/openapi/openapi.yaml
```

The `PersonRepository.go` file acts as a central point for connecting the created database and api server. The API server has the following calls:

**1. Person Repository**

- `GET api/v1/persons` - gets all persons in database
- `POST api/v1/persons` - creates new person in database
- `DELETE api/v1/persons` - truncates database persons table
- `GET api/v1/persons/summary` - gets shortened list of person data
- `GET api/v1/persons/{id}` - get specific person by ID
- `PUT api/v1/persons/{id}` - update person by ID
- `DELETE api/v1/persons/{id}` - delete person by ID

**2. Database Management**

- `GET api/v1/backups` - get all backups on device
- `POST api/v1/backups` - create new backup file
- `DELETE api/v1/backups/{filename}` - delete backup file
- `POST api/v1/backups/{filename}/restore` - restore backup

### User Interface

TODO

## API Application Deployment

### Prerequisites

The application requires installing the following tools:
- `go`,
- `gcc`.

Additionally, `go-sqlite3` requires `CGO_ENABLED` to work. Set the flag to true before building and running the project.

```bash
export CGO_ENABLED=1
```

### Build and Run

To build the application, run the code below.

```bash
mkdir bin && cd bin
go build -o pop-db ../cmd/
```

To run the application (on `http://localhost:8080/swagger/`), run the following code:

```bash
./pop-db
```

## UI Application Deployment

TODO

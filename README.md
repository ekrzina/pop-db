# PopDB

## Introduction

PopDB is a lightweight, personal-use database project designed to manage and display population census data through a dynamic web application. The project integrates a relational SQL database with a streamlined frontend and backend architecture.

## Technology Stack

- `SQLLite` - embedded relational database for data storage
- `Next.js + Tailwind + shadcn/ui + Tanstack Query` - lightweight frontend interface for client-side interactions
- `Go` - backend server-side logic using the Gin web framework

## Prerequisites

Before you can build or run the project you need to have the following
tools installed on your development machine (and on any machine where you
intend to build or package the application):

- **Go 1.21+** (with `gcc`/`clang` available for `cgo`; the `go-sqlite3` driver
  requires `CGO_ENABLED=1`)
- **Node.js 18+** (or compatible; used to build and run the Next.js frontend)
- **npm** (bundled with Node.js) or **pnpm/yarn** as an alternative package
  manager
- **git** (for cloning the repository, optional once source is local)
- **Bash-compatible shell** (for using the provided `start.sh`/`install.sh`
  scripts; Linux/macOS have one built‑in, Windows users can use Git Bash, WSL,
  Cygwin, etc.)
- On Windows, if you don’t have a shell you can still run the bundled `dist`
  artefacts directly – see the packaging instructions above.

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

The PopDB user interface is a web-based frontend built with Next.js, styled using Tailwind CSS and components from shadcn/ui, and uses TanStack Query for efficient client-side data fetching. It communicates with the backend API to display and manage population data in real time.

**Key features of the UI**:
- Person List: Scrollable buttons representing all persons in the database. Selecting a person displays detailed information in a card.
- Search & Filtering: Users can search persons by name, surname, or occupation.
- Pagination & Navigation: Keyboard arrow keys or on-screen buttons navigate between persons and pages.
- Create & Edit Modals: Pop-up forms for adding or updating a person, with validation to ensure required fields (name, surname, city, nationality) and medical data (height and weight > 0) are correctly filled.
- Delete Functionality: Users can remove a person with confirmation prompts to prevent accidental deletions.
- Summary Page: Quick overview of stored persons and key statistics.
- Responsive Design: Works across desktop and tablet screens, with smooth scrolling and accessible buttons.

All UI interactions trigger API calls to the backend, ensuring that the frontend state reflects the latest database content. Optimistic updates and loading indicators are used to improve the user experience during network requests.

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

To run the UI application in development mode use:

```bash
npm run popdb-ui/dev
```

`favicon` phoenix logo downloaded here:
```html
<a href="" title="phoenix icons">Phoenix icons created by Freepik - Flaticon</a>
```

### Building for production

```bash
cd popdb-ui
npm install      # only the first time or when dependencies change
npm run build     # creates the `.next` directory with production assets
npm run start     # serve the built app on port 3000
```

### Combined startup script

A start script `start.sh` has been created for building and running purposes.

```bash
chmod +x start.sh   # change permissions if needed
```

To **build and run** on your device use the following:

```bash
./start.sh
```

To **build and package** the app for a certain distribution, start the script like so:

```bash
./start.sh package
```

This will package the application frontend and backend builds for your chosen distribution.

The file structure should be the following:

```bash
dist/
├─ backend/
│  ├─ popdb
│  ├─ config.yaml
│  ├─ dbase/
│  └─ backup/
└─ frontend/
   ├─ .next/
   ├─ public/
   ├─ package.json
   └─ node_modules/
```

Use the `run.sh` script to run distribution packages.

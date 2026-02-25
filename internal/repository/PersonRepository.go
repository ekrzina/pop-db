package repository

import (
	"pop-db/internal/db"
	"pop-db/internal/repository/models"
)

// PersonRepository stores database object for access to SQL database
type PersonRepository struct {
	db *db.SQLLiteDB
}

// CreatePerson creates a single person and adds it to database file
// Parameters:
//   - p: New Person to add to database.
//
// Returns:
//   - int64: Last inserted person ID (auto-incremented).
//   - error: Error on insertion fail.
func (r *PersonRepository) CreatePerson(p *models.Person) (int64, error) {
	result, err := r.db.Execute(`
		INSERT INTO person (name, surname, occupation, date_of_birth, nationality, city, notes, picture)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Name,
		p.Surname,
		p.Occupation,
		p.DateOfBirth,
		p.Nationality,
		p.City,
		p.Notes,
		p.Picture,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetPersonByID fetches a single person that matches selected identification number of person
// Parameters:
//   - id: Identification number of person to fetch.
//
// Returns:
//   - *models.Person: Person information fetched from database if ID matches.
//   - error: Error on search fail.
func (r *PersonRepository) GetPersonByID(id int64) (*models.Person, error) {
	row := r.db.QueryRow(`
		SELECT id, name, surname, occupation, date_of_birth, nationality, city, notes, picture
		FROM person WHERE id = ?`, id)
	var p models.Person
	err := row.Scan(
		&p.ID,
		&p.Name,
		&p.Surname,
		&p.Occupation,
		&p.DateOfBirth,
		&p.Nationality,
		&p.City,
		&p.Notes,
		&p.Picture,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ListPersons lists all IDs, names and surnames of persons stored in database.
// Returns:
//   - []models.Person: List of all Person information fetched from database.
//   - error: Error on query fail.
func (r *PersonRepository) ListPersons() ([]models.Person, error) {
	rows, err := r.db.Query(`
		SELECT id, name, surname FROM person`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Scan all persons from database and append each to list
	var persons []models.Person
	for rows.Next() {
		var p models.Person
		if err := rows.Scan(&p.ID, &p.Name, &p.Surname); err != nil {
			return nil, err
		}
		persons = append(persons, p)
	}
	return persons, nil
}

// DeletePerson deletes a single person from database whose identification number matches selected ID
// Parameters:
//   - id: Identification number of person to delete from database.
//
// Returns:
//
//	error: Error on query fail.
func (r *PersonRepository) DeletePerson(id int64) error {
	_, err := r.db.Execute(`DELETE FROM person WHERE id = ?`, id)
	return err
}

// UpdatePerson updates a single person from database whose identification number matches selected ID
// Parameters:
//   - p: Person object to update data for.
//
// Returns:
//
//	error: Error on query fail.
func (r *PersonRepository) UpdatePerson(p *models.Person) error {
	_, err := r.db.Execute(`
		UPDATE person
		SET name = ?, surname = ?, occupation = ?, date_of_birth = ?, nationality = ?, city = ?, notes = ?, picture = ?
		WHERE id = ?`,
		p.Name,
		p.Surname,
		p.Occupation,
		p.DateOfBirth,
		p.Nationality,
		p.City,
		p.Notes,
		p.Picture,
		p.ID,
	)
	return err
}

// CreateMedicalData creates a single person's medical information and adds it to database file
// Parameters:
//   - m: New MedicalData to add to person in database.
//
// Returns:
//   - error: Error on insertion fail.
func (r *PersonRepository) CreateMedicalData(m *models.MedicalData) error {
	_, err := r.db.Execute(`
		INSERT INTO medical_data (person_id, height, weight, blood_type, medical_conditions)
		VALUES (?, ?, ?, ?, ?)`,
		m.PersonID,
		m.Height,
		m.Weight,
		m.BloodType,
		m.MedicalConditions,
	)
	return err
}

// GetPersonWithMedicalData fetches a single person's full data consisting of personal and medical information
// Parameters:
//   - m: New MedicalData to add to person in database.
//
// Returns:
//   - error: Error on insertion fail.
func (r *PersonRepository) GetPersonWithMedicalData(id int64) (*models.Person, error) {
	row := r.db.QueryRow(`
		SELECT p.id, p.name, p.surname, p.occupation, p.date_of_birth, p.nationality, p.city, p.notes, p.picture,
		       m.height, m.weight, m.blood_type, m.medical_conditions
		FROM person p
		LEFT JOIN medical_data m ON p.id = m.person_id
		WHERE p.id = ?`, id)

	var p models.Person
	var m models.MedicalData

	err := row.Scan(
		&p.ID,
		&p.Name,
		&p.Surname,
		&p.Occupation,
		&p.DateOfBirth,
		&p.Nationality,
		&p.City,
		&p.Notes,
		&p.Picture,
		&m.Height,
		&m.Weight,
		&m.BloodType,
		&m.MedicalConditions,
	)
	if err != nil {
		return nil, err
	}

	m.PersonID = p.ID
	p.Medical = &m

	return &p, nil
}

// CreateFullPerson creates a single person with medical data and adds it to database file with transaction support
// Parameters:
//   - p: New Person to add to database.
//
// Returns:
//   - int64: Last inserted person ID (auto-incremented).
//   - error: Error on insertion fail.
func (r *PersonRepository) CreateFullPerson(
	p *models.Person,
	m *models.MedicalData,
) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	result, err := tx.Exec(`
		INSERT INTO person (name, surname, occupation, date_of_birth, nationality, city, notes, picture)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Name,
		p.Surname,
		p.Occupation,
		p.DateOfBirth,
		p.Nationality,
		p.City,
		p.Notes,
		p.Picture,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if m != nil {
		_, err = tx.Exec(`
			INSERT INTO medical_data (person_id, height, weight, blood_type, medical_conditions)
			VALUES (?, ?, ?, ?, ?)`,
			id,
			m.Height,
			m.Weight,
			m.BloodType,
			m.MedicalConditions,
		)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

// NewPersonRepository returns PersonRepository object
func NewPersonRepository(db *db.SQLLiteDB) *PersonRepository {
	return &PersonRepository{db: db}
}

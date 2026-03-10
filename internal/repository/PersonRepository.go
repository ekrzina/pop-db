package repository

import (
	"database/sql"

	"github.com/haoli/pop-db/internal/dbman"
	"github.com/haoli/pop-db/internal/repository/models"

	"github.com/rs/zerolog"
)

// PersonRepository stores database object for access to SQL database
type PersonRepository struct {
	manager *dbman.DbManager
	logger  *zerolog.Logger
}

// GetManager returns the database manager instance for database operations
// Returns:
//   - *DbManager: The database manager instance.
func (r *PersonRepository) Manager() *dbman.DbManager {
	return r.manager
}
func (r *PersonRepository) GetManager() *dbman.DbManager {
	return r.manager
}

// safeRollback safely executes the rollback function on transaction fail
// Parameters:
//   - tx: SQL TX data.
func (r *PersonRepository) safeRollback(tx *sql.Tx) {
	if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
		r.logger.Error().Err(err).Msg("Rollback failed")
	}
}

// CreatePerson creates a single person and adds it to database file
// Parameters:
//   - p: New Person to add to database.
//
// Returns:
//   - int64: Last inserted person ID (auto-incremented).
//   - error: Error on insertion fail.
func (r *PersonRepository) CreatePerson(p *models.Person) (int64, error) {
	result, err := r.manager.DB.Execute(`
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

// CreateFullPerson creates a person and optional medical data in a single transaction
// Parameters:
//   - p: New Person to add to database.
//   - m: Optional MedicalData to add to database.
//
// Returns:
//   - int64: Last inserted person ID (auto-incremented).
//   - error: Error on insertion fail.
func (r *PersonRepository) CreateFullPerson(p *models.Person, m *models.MedicalData) (int64, error) {
	tx, err := r.manager.DB.Begin()
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
		r.safeRollback(tx)
		return 0, err
	}
	personID, err := result.LastInsertId()
	if err != nil {
		r.safeRollback(tx)
		return 0, err
	}
	if m != nil {
		_, err = tx.Exec(`
            INSERT INTO medical_data (person_id, height, weight, blood_type, medical_conditions)
            VALUES (?, ?, ?, ?, ?)`,
			personID,
			m.Height,
			m.Weight,
			m.BloodType,
			m.MedicalConditions,
		)
		if err != nil {
			r.safeRollback(tx)
			return 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return personID, nil
}

// GetPersonByID fetches a single person that matches selected identification number of person
// Parameters:
//   - id: Identification number of person to fetch.
//
// Returns:
//   - *models.Person: Person information fetched from database if ID matches.
//   - error: Error on search fail.
func (r *PersonRepository) GetPersonByID(id int64) (*models.Person, error) {
	row := r.manager.DB.QueryRow(`
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
//   - []models.PersonSummary: List of all PersonSummary information fetched from database.
//   - error: Error on query fail.
func (r *PersonRepository) ListPersons() ([]models.PersonSummary, error) {
	rows, err := r.manager.DB.Query(`
		SELECT id, name, surname FROM person`)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			r.logger.Error().Err(err).Msg("Failed to close rows")
		}
	}()
	// Scan all persons from database and append each to list
	var persons []models.PersonSummary
	for rows.Next() {
		var p models.PersonSummary
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
//   - error: Error on query fail.
func (r *PersonRepository) DeletePerson(id int64) error {
	_, err := r.manager.DB.Execute(`DELETE FROM person WHERE id = ?`, id)
	return err
}

// TruncatePersons deletes all persons from database
// Returns:
//   - int: Number of rows affected by deletion.
//   - error: Error on query fail.
func (r *PersonRepository) TruncatePersons() (int, error) {
	res, err := r.manager.DB.Execute(`DELETE FROM person`)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return int(rows), err
}

// UpdateFullPerson updates a single person from database whose identification number matches selected ID
// Parameters:
//   - p: Person object to update data for.
//   - m: MedicalData object to update data for.
//
// Returns:
//
//	error: Error on query fail.
func (r *PersonRepository) UpdateFullPerson(p *models.Person, m *models.MedicalData) error {
	tx, err := r.manager.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
        UPDATE person
        SET name = ?, surname = ?, occupation = ?, date_of_birth = ?, nationality = ?, city = ?, notes = ?, picture = ?
        WHERE id = ?`,
		p.Name, p.Surname, p.Occupation, p.DateOfBirth, p.Nationality, p.City, p.Notes, p.Picture, p.ID,
	)
	if err != nil {
		r.safeRollback(tx)
		return err
	}
	if m != nil {
		res, err := tx.Exec(`
        UPDATE medical_data
        SET height = ?, weight = ?, blood_type = ?, medical_conditions = ?
        WHERE person_id = ?`,
			m.Height, m.Weight, m.BloodType, m.MedicalConditions, m.PersonID,
		)
		if err != nil {
			r.safeRollback(tx)
			return err
		}
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			_, err = tx.Exec(`
            INSERT INTO medical_data (person_id, height, weight, blood_type, medical_conditions)
            VALUES (?, ?, ?, ?, ?)`,
				m.PersonID, m.Height, m.Weight, m.BloodType, m.MedicalConditions,
			)
			if err != nil {
				r.safeRollback(tx)
				return err
			}
		}
	}
	return tx.Commit()
}

// CreateMedicalData creates a single person's medical information and adds it to database file
// Parameters:
//   - m: New MedicalData to add to person in database.
//
// Returns:
//   - error: Error on insertion fail.
func (r *PersonRepository) CreateMedicalData(m *models.MedicalData) error {
	_, err := r.manager.DB.Execute(`
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

// GetPersonsWithMedicalData returns all persons with their medical data
// Returns:
//   - []models.Person: A list of all person models.
//   - error: Error returned on reading fail.
func (r *PersonRepository) GetPersonsWithMedicalData() ([]models.Person, error) {
	rows, err := r.manager.DB.Query(`
        SELECT 
            p.id, p.name, p.surname, p.occupation, p.date_of_birth, 
            p.nationality, p.city, p.notes, p.picture,
            m.height, m.weight, m.blood_type, m.medical_conditions
        FROM person p
        LEFT JOIN medical_data m ON p.id = m.person_id
    `)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			r.logger.Error().Err(err).Msg("Failed to close rows")
		}
	}()
	var persons []models.Person

	for rows.Next() {
		var p models.Person
		// Nullable medical fields
		var (
			height     sql.NullFloat64
			weight     sql.NullFloat64
			bloodType  sql.NullString
			conditions sql.NullString
		)
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Surname,
			&p.Occupation,
			&p.DateOfBirth,
			&p.Nationality,
			&p.City,
			&p.Notes,
			&p.Picture,
			&height,
			&weight,
			&bloodType,
			&conditions,
		)
		if err != nil {
			return nil, err
		}
		if height.Valid || weight.Valid || bloodType.Valid || conditions.Valid {
			md := models.MedicalData{
				Height:    height.Float64,
				Weight:    weight.Float64,
				BloodType: bloodType.String,
			}
			if conditions.Valid {
				md.MedicalConditions = conditions.String
			}
			p.Medical = &md
		}
		persons = append(persons, p)
	}
	return persons, nil
}

// GetPersonWithMedicalData fetches a single person's full data consisting of personal and medical information
// Parameters:
//   - id: Identification number of person with medical data.
//
// Returns:
//   - *models.Person: Person model that matches identification number.
//   - error: Error on reading or id matching fail.
func (r *PersonRepository) GetPersonWithMedicalData(id int64) (*models.Person, error) {
	row := r.manager.DB.QueryRow(`
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

// NewPersonRepository returns PersonRepository object
// Parameters:
//   - *dbman.DbManager: Database manager for person repository.
//   - logger: Zerolog logger for person repository.
//
// Returns:
//   - *PersonRepository: Person repository to manage persons with.
func NewPersonRepository(manager *dbman.DbManager, logger *zerolog.Logger) *PersonRepository {
	return &PersonRepository{
		manager: manager,
		logger:  logger,
	}
}

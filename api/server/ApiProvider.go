package server

import (
	"context"
	"net/http"
	"time"

	"github.com/haoli/pop-db/internal/repository"
	"github.com/haoli/pop-db/internal/repository/models"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/rs/zerolog"
)

type ApiProvider struct {
	logger *zerolog.Logger
	repo   repository.PersonRepository
}

// GET api/v1/persons
func (a *ApiProvider) GetApiV1Persons(ctx context.Context, request GetApiV1PersonsRequestObject) (GetApiV1PersonsResponseObject, error) {
	persons, err := a.repo.GetPersonsWithMedicalData()
	if err != nil {
		return GetApiV1Persons500JSONResponse{
			InternalServerErrorJSONResponse{
				Error:     "internal_error",
				Message:   "Failed to fetch persons: " + err.Error(),
				Status:    http.StatusInternalServerError,
				Timestamp: time.Now(),
			},
		}, nil
	}
	response := make([]Person, 0, len(persons))
	for _, p := range persons {
		var medical *MedicalData
		if p.Medical != nil {
			medical = &MedicalData{
				BloodType:         MedicalDataBloodType(p.Medical.BloodType),
				Height:            float32(p.Medical.Height),
				Weight:            float32(p.Medical.Weight),
				MedicalConditions: &p.Medical.MedicalConditions,
			}
		}
		response = append(response, Person{
			Id:          int(p.ID),
			Name:        p.Name,
			Surname:     p.Surname,
			City:        p.City,
			Nationality: p.Nationality,
			DateOfBirth: openapi_types.Date{Time: p.DateOfBirth},
			Occupation:  &p.Occupation,
			Notes:       &p.Notes,
			Medical:     medical,
		})
	}
	return GetApiV1Persons200JSONResponse(response), nil
}

// POST api/v1/persons
func (a *ApiProvider) PostApiV1Persons(ctx context.Context, request PostApiV1PersonsRequestObject) (PostApiV1PersonsResponseObject, error) {
	r := request.Body
	p := &models.Person{
		Name:        r.Name,
		Surname:     r.Surname,
		Occupation:  *r.Occupation,
		Nationality: r.Nationality,
		City:        r.City,
		Notes:       *r.Notes,
		DateOfBirth: r.DateOfBirth.Time,
	}
	var conditions string
	if r.Medical.MedicalConditions != nil {
		conditions = *r.Medical.MedicalConditions
	}
	m := &models.MedicalData{
		Height:            float64(r.Medical.Height),
		Weight:            float64(r.Medical.Weight),
		BloodType:         string(r.Medical.BloodType),
		MedicalConditions: conditions,
	}
	_, err := a.repo.CreateFullPerson(p, m)
	if err != nil {
		return PostApiV1Persons500JSONResponse{
			InternalServerErrorJSONResponse{
				Error:     "internal_error",
				Message:   "Failed to create person: " + err.Error(),
				Status:    http.StatusInternalServerError,
				Timestamp: time.Now(),
			},
		}, nil
	}
	return PostApiV1Persons204Response{}, nil
}

// DELETE api/v1/persons
func (a *ApiProvider) DeleteApiV1Persons(ctx context.Context, request DeleteApiV1PersonsRequestObject) (DeleteApiV1PersonsResponseObject, error) {
	rows, err := a.repo.TruncatePersons()
	if err != nil {
		return DeleteApiV1Persons500JSONResponse{
			InternalServerErrorJSONResponse{
				Error:     "internal_error",
				Message:   "Failed to delete persons: " + err.Error(),
				Status:    http.StatusNotFound,
				Timestamp: time.Now(),
			},
		}, nil
	}
	msg := "All persons deleted successfully"
	return DeleteApiV1Persons200JSONResponse{
		Message:      &msg,
		DeletedCount: &rows,
	}, nil
}

// GET api/v1/persons/summary
func (a *ApiProvider) GetApiV1PersonsSummary(ctx context.Context, request GetApiV1PersonsSummaryRequestObject) (GetApiV1PersonsSummaryResponseObject, error) {
	persons, err := a.repo.ListPersons()
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to fetch persons")
		return GetApiV1PersonsSummary500JSONResponse{
			InternalServerErrorJSONResponse{
				Error:     "internal_error",
				Message:   "Failed to fetch persons: " + err.Error(),
				Status:    http.StatusInternalServerError,
				Timestamp: time.Now(),
			},
		}, nil
	}
	response := make([]PersonSummary, 0, len(persons))
	for _, p := range persons {
		response = append(response, PersonSummary{
			Id:      int(p.ID),
			Name:    p.Name,
			Surname: p.Surname,
		})
	}
	return GetApiV1PersonsSummary200JSONResponse(response), nil
}

// GET api/v1/persons/{id}
func (a *ApiProvider) GetApiV1PersonsId(ctx context.Context, request GetApiV1PersonsIdRequestObject) (GetApiV1PersonsIdResponseObject, error) {
	p, err := a.repo.GetPersonWithMedicalData(int64(request.Id))
	if err != nil {
		return GetApiV1PersonsId404JSONResponse{
			NotFoundJSONResponse{
				Error:     "not_found",
				Message:   "Person not found: " + err.Error(),
				Status:    http.StatusNotFound,
				Timestamp: time.Now(),
			},
		}, nil
	}
	response := Person{
		Id:          int(p.ID),
		Name:        p.Name,
		Surname:     p.Surname,
		City:        p.City,
		Nationality: p.Nationality,
		DateOfBirth: openapi_types.Date{Time: p.DateOfBirth},
		Occupation:  &p.Occupation,
		Notes:       &p.Notes,
		Medical: &MedicalData{
			BloodType:         MedicalDataBloodType(p.Medical.BloodType),
			Height:            float32(p.Medical.Height),
			Weight:            float32(p.Medical.Weight),
			MedicalConditions: &p.Medical.MedicalConditions,
		},
	}
	return GetApiV1PersonsId200JSONResponse(response), nil
}

// PUT api/v1/persons/{id}
func (a *ApiProvider) PutApiV1PersonsId(ctx context.Context, request PutApiV1PersonsIdRequestObject) (PutApiV1PersonsIdResponseObject, error) {
	r := request.Body
	p := &models.Person{
		ID:          int64(r.Id),
		Name:        r.Name,
		Surname:     r.Surname,
		Occupation:  *r.Occupation,
		Nationality: r.Nationality,
		City:        r.City,
		Notes:       *r.Notes,
		DateOfBirth: r.DateOfBirth.Time,
	}
	var m *models.MedicalData
	if r.Medical != nil {
		m = &models.MedicalData{
			PersonID:          int64(r.Id),
			Height:            float64(r.Medical.Height),
			Weight:            float64(r.Medical.Weight),
			BloodType:         string(r.Medical.BloodType),
			MedicalConditions: *r.Medical.MedicalConditions,
		}
	}
	err := a.repo.UpdateFullPerson(p, m)
	if err != nil {
		return PutApiV1PersonsId400JSONResponse{
			BadRequestJSONResponse{
				Error:     "bad_request",
				Message:   "Person update failed: " + err.Error(),
				Status:    http.StatusBadRequest,
				Timestamp: time.Now(),
			},
		}, nil
	}
	return PutApiV1PersonsId204Response{}, nil
}

// DELETE api/v1/persons/{id}
func (a *ApiProvider) DeleteApiV1PersonsId(ctx context.Context, request DeleteApiV1PersonsIdRequestObject) (DeleteApiV1PersonsIdResponseObject, error) {
	err := a.repo.DeletePerson(int64(request.Id))
	if err != nil {
		return DeleteApiV1PersonsId404JSONResponse{
			NotFoundJSONResponse{
				Error:     "not_found",
				Message:   "Person not found: " + err.Error(),
				Status:    http.StatusNotFound,
				Timestamp: time.Now(),
			},
		}, nil
	}
	return DeleteApiV1PersonsId204Response{}, nil
}

// GET api/v1/backups
func (a *ApiProvider) GetApiV1Backups(ctx context.Context, request GetApiV1BackupsRequestObject) (GetApiV1BackupsResponseObject, error) {
	man := a.repo.Manager()
	backups, err := man.ListBackups()
	if err != nil {
		return GetApiV1Backups500JSONResponse{
			InternalServerErrorJSONResponse{
				Error:     "internal_server_error",
				Message:   "Failed to list backups: " + err.Error(),
				Status:    http.StatusInternalServerError,
				Timestamp: time.Now(),
			},
		}, nil
	}
	response := make([]Backup, len(backups))
	for i, b := range backups {
		response[i] = Backup{
			Filename:  b.Filename,
			Path:      b.Path,
			CreatedAt: b.CreatedAt,
			SizeBytes: int(b.SizeBytes),
		}
	}
	return GetApiV1Backups200JSONResponse(response), nil
}

// POST api/v1/backups
func (a *ApiProvider) PostApiV1Backups(ctx context.Context, request PostApiV1BackupsRequestObject) (PostApiV1BackupsResponseObject, error) {
	man := a.repo.Manager()
	meta, err := man.WriteBackup()
	if err != nil {
		return PostApiV1Backups500JSONResponse{
			InternalServerErrorJSONResponse{
				Error:     "internal_server_error",
				Message:   "Failed to add backup: " + err.Error(),
				Status:    http.StatusInternalServerError,
				Timestamp: time.Now(),
			},
		}, nil
	}
	backup := Backup{
		Filename:  meta.Filename,
		Path:      meta.Path,
		CreatedAt: meta.CreatedAt,
		SizeBytes: int(meta.SizeBytes),
	}
	return PostApiV1Backups200JSONResponse(backup), nil
}

// DELETE api/v1/backups/{filename}
func (a *ApiProvider) DeleteApiV1BackupsFilename(ctx context.Context, request DeleteApiV1BackupsFilenameRequestObject) (DeleteApiV1BackupsFilenameResponseObject, error) {
	man := a.repo.Manager()
	err := man.DeleteBackup(request.Filename)
	if err != nil {
		return DeleteApiV1BackupsFilename404JSONResponse{
			NotFoundJSONResponse{
				Error:     "not_found",
				Message:   "Backup not found: " + err.Error(),
				Status:    http.StatusNotFound,
				Timestamp: time.Now(),
			},
		}, nil
	}
	return DeleteApiV1BackupsFilename204Response{}, nil
}

// POST api/v1/backups/{filename}/restore
func (a *ApiProvider) PostApiV1BackupsFilenameRestore(ctx context.Context, request PostApiV1BackupsFilenameRestoreRequestObject) (PostApiV1BackupsFilenameRestoreResponseObject, error) {
	man := a.repo.Manager()
	err := man.RestoreBackup(request.Filename)
	if err != nil {
		return PostApiV1BackupsFilenameRestore404JSONResponse{
			NotFoundJSONResponse{
				Error:     "not_found",
				Message:   "Backup not found: " + err.Error(),
				Status:    http.StatusNotFound,
				Timestamp: time.Now(),
			},
		}, nil
	}
	resp := RestoreResponse{
		Message:      "Database restored successfully.",
		RestoredFrom: request.Filename,
	}
	return PostApiV1BackupsFilenameRestore200JSONResponse(resp), nil
}

func NewApiProvider(logger *zerolog.Logger, repo repository.PersonRepository) *ApiProvider {
	return &ApiProvider{
		logger: logger,
		repo:   repo,
	}
}

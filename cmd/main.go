package main

import (
	"bytes"
	"embed"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"pop-db/internal/db"
	"pop-db/internal/repository"
	"pop-db/internal/repository/models"
)

//go:embed configs/config.yaml
var configFS embed.FS

func initializeLogger() zerolog.Logger {
	// Console writer for human-readable logs
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	// Global time format for structured logs
	zerolog.TimeFieldFormat = time.RFC3339

	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()
	// Replace the standard logger
	log.Logger = logger
	return logger
}

func loadConfig() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	configData, err := configFS.ReadFile("configs/config.yaml")
	if err != nil {
		return nil, err
	}

	if err := v.ReadConfig(bytes.NewBuffer(configData)); err != nil {
		return nil, err
	}

	return v, nil
}

func main() {
	logger := initializeLogger()

	v, err := loadConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load config")
	}

	// Ensure directories exist
	os.MkdirAll(v.GetString("database.path"), os.ModePerm)
	os.MkdirAll(v.GetString("database.backupPath"), os.ModePerm)

	database, err := db.NewSQLiteDB(v, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer database.Close()

	repo := repository.NewPersonRepository(database)

	// Example data
	person := &models.Person{
		Name:        "Arthur",
		Surname:     "Dent",
		Occupation:  "Human",
		DateOfBirth: time.Date(1978, 3, 11, 0, 0, 0, 0, time.UTC),
		Nationality: "British",
		City:        "London",
		Notes:       "Mostly harmless.",
	}

	medical := &models.MedicalData{
		Height:            178,
		Weight:            75,
		BloodType:         "O+",
		MedicalConditions: "Existential confusion.",
	}

	id, err := repo.CreateFullPerson(person, medical)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to insert person")
	}
	logger.Info().Int64("id", id).Msg("Inserted person")

	fullPerson, err := repo.GetPersonWithMedicalData(id)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to fetch person")
	}

	logger.Info().Interface("person", fullPerson).Msg("Fetched person")

	if err := database.WriteBackup(); err != nil {
		logger.Error().Err(err).Msg("Backup failed")
	} else {
		logger.Info().Msg("Backup completed")
	}
}

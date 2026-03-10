package models

import "time"

// Structure to contain information on a single person
type Person struct {
	ID          int64
	Name        string
	Surname     string
	Occupation  string
	DateOfBirth time.Time
	Nationality string
	City        string
	Notes       string
	Picture     []byte
	Medical     *MedicalData
}

// Structure to contain a single person's medical / identification data representation
type MedicalData struct {
	PersonID          int64
	Height            float64
	Weight            float64
	BloodType         string
	MedicalConditions string
}

// Structure to contain a single person's summary data representation
type PersonSummary struct {
	ID      int64
	Name    string
	Surname string
}

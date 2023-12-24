package domain

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Login       string
	Password    string
	IsActivated bool
	Employee    uuid.UUID
}

type Department struct {
	ID   uuid.UUID
	Name string
}

type Employee struct {
	ID          uuid.UUID
	Surname     string
	Name        string
	Patronymic  string
	PhoneNumber string
	Password    string
	Position    uuid.UUID
	Department  uuid.UUID
}

type Position struct {
	ID               uuid.UUID
	Title            string
	MaxPerDepartment uint
}

type ProfessionalDevelopment struct {
	ID           uuid.UUID
	Employee     uuid.UUID
	Title        string
	StartingDate time.Time
	EndingDate   time.Time
}

type Incident struct {
	ID          uuid.UUID
	Employee    uuid.UUID
	Description string
	Date        time.Time
}

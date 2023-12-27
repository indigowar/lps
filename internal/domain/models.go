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
	ID          uuid.UUID `db:"id"`
	Surname     string    `db:"surname"`
	Name        string    `db:"name"`
	Patronymic  *string   `db:"patronymic"`
	PhoneNumber string    `db:"phone_number"`
	Position    uuid.UUID `db:"position"`
	Department  uuid.UUID `db:"department"`
}

type Position struct {
	ID    uuid.UUID
	Title string
	Level string
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

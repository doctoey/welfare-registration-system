package entities

import "time"

type Application struct {
	ID              int       `json:"id"`
	ReferenceNumber string    `json:"referenceNumber"`
	CitizenID       string    `json:"citizenId"`
	FullName        string    `json:"fullName"`
	BirthDate       string    `json:"birthDate"`
	AnnualIncome    float64   `json:"annualIncome"`
	CurrentAddress  string    `json:"currentAddress"`
	Status          string    `json:"status"` // "pending", "approved", "rejected"
	Reason          string    `json:"reason,omitempty"`
	SubmittedAt     time.Time `json:"submittedAt"`
}

type CreateApplicationDTO struct {
	CitizenID      string  `json:"citizenId"`
	FullName       string  `json:"fullName"`
	BirthDate      string  `json:"birthDate"`
	AnnualIncome   float64 `json:"annualIncome"`
	CurrentAddress string  `json:"currentAddress"`
}

type StatusResponseDTO struct {
	ReferenceNumber string `json:"referenceNumber"`
	CitizenID       string `json:"citizenId"`
	FullName        string `json:"fullName"`
	Status          string `json:"status"`
	Reason          string `json:"reason,omitempty"`
}

type UpdateStatusDTO struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type ApplicationRepository interface {
	FindAll(statusFilter string) ([]Application, error)
	FindByCitizenID(citizenID string) (*Application, error)
	FindByID(id int) (*Application, error)
	Save(app *Application) error
	Update(app *Application) error
}

type ApplicationUsecase interface {
	Register(dto CreateApplicationDTO) (string, error)
	GetStatus(citizenID string, birthDate string, refNumber string) (*StatusResponseDTO, error)
	GetOfficerApplications(statusFilter string) ([]Application, error)
	ReviewApplication(id int, dto UpdateStatusDTO) (*Application, error)
}

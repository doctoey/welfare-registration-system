package entities

import "time"

type Application struct {
	ID              string    `json:"id" gorm:"type:varchar(36);primaryKey"`
	ReferenceNumber string    `json:"referenceNumber" gorm:"type:varchar(30);uniqueIndex"`
	CitizenID       string    `json:"citizenId" gorm:"type:varchar(13);uniqueIndex;not null"`
	FullName        string    `json:"fullName" gorm:"type:varchar(200);not null"`
	FirstName       string    `json:"firstName,omitempty" gorm:"type:varchar(100)"`
	LastName        string    `json:"lastName,omitempty" gorm:"type:varchar(100)"`
	BirthDate       string    `json:"birthDate" gorm:"type:date;not null"`
	PhoneNumber     string    `json:"phoneNumber,omitempty" gorm:"type:varchar(15)"`
	AnnualIncome    float64   `json:"annualIncome" gorm:"type:decimal(12,2);not null"`
	CurrentAddress  string    `json:"currentAddress" gorm:"column:address_detail;type:text;not null"`
	Status          string    `json:"status" gorm:"type:varchar(20);index;not null"` // "pending", "approved", "rejected"
	Reason          string    `json:"reason,omitempty" gorm:"column:rejection_reason;type:text"`
	CreatedAt       time.Time `json:"submittedAt" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (Application) TableName() string {
	return "applications"
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
	FindByID(id string) (*Application, error)
	Save(app *Application) error
	Update(app *Application) error
}

type ApplicationUsecase interface {
	Register(dto CreateApplicationDTO) (string, error)
	GetStatus(citizenID string, birthDate string, refNumber string) (*StatusResponseDTO, error)
	GetOfficerApplications(statusFilter string) ([]Application, error)
	ReviewApplication(id string, dto UpdateStatusDTO) (*Application, error)
}

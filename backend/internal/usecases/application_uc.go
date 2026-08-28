package usecases

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"
	"welfare-registration-backend/internal/entities"

	"github.com/google/uuid"
)

type applicationUsecase struct {
	repo entities.ApplicationRepository
}

func NewApplicationUsecase(repo entities.ApplicationRepository) entities.ApplicationUsecase {
	return &applicationUsecase{repo: repo}
}

func isValidThaiID(id string) bool {
	clean := strings.ReplaceAll(id, "-", "")
	if len(clean) != 13 {
		return false
	}
	sum := 0
	for i := 0; i < 12; i++ {
		digit := int(clean[i] - '0')
		sum += digit * (13 - i)
	}
	checkDigit := (11 - (sum % 11)) % 10
	return checkDigit == int(clean[12]-'0')
}

func maskName(name string) string {
	parts := strings.Fields(name)
	var maskedParts []string
	for _, p := range parts {
		runes := []rune(p)
		if len(runes) <= 2 {
			maskedParts = append(maskedParts, string(runes))
		} else {
			maskedParts = append(maskedParts, string(runes[:2])+strings.Repeat("*", utf8.RuneCountInString(p)-2))
		}
	}
	return strings.Join(maskedParts, " ")
}

func maskCitizenID(id string) string {
	clean := strings.ReplaceAll(id, "-", "")
	if len(clean) != 13 {
		return id
	}
	return fmt.Sprintf("%s-%s-*****-**-%s", clean[0:1], clean[1:5], clean[12:13])
}

func (u *applicationUsecase) Register(dto entities.CreateApplicationDTO) (string, error) {
	cleanCitizenID := strings.TrimSpace(dto.CitizenID)
	if cleanCitizenID == "" || strings.TrimSpace(dto.FullName) == "" || dto.BirthDate == "" {
		return "", errors.New("REQUIRED_FIELDS_MISSING")
	}

	if !isValidThaiID(cleanCitizenID) {
		return "", errors.New("INVALID_CITIZEN_ID_CHECKSUM")
	}

	existing, _ := u.repo.FindByCitizenID(cleanCitizenID)
	if existing != nil {
		return "", errors.New("CITIZEN_ALREADY_REGISTERED")
	}

	names := strings.Fields(strings.TrimSpace(dto.FullName))
	firstName := ""
	lastName := ""
	if len(names) > 0 {
		firstName = names[0]
	}
	if len(names) > 1 {
		lastName = strings.Join(names[1:], " ")
	}

	refNumber := fmt.Sprintf("WRS-2026-%06d", rand.Intn(900000)+100000)
	newApp := entities.Application{
		ID:              uuid.New().String(),
		ReferenceNumber: refNumber,
		CitizenID:       cleanCitizenID,
		FullName:        strings.TrimSpace(dto.FullName),
		FirstName:       firstName,
		LastName:        lastName,
		BirthDate:       dto.BirthDate,
		AnnualIncome:    dto.AnnualIncome,
		CurrentAddress:  strings.TrimSpace(dto.CurrentAddress),
		Status:          "pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := u.repo.Save(&newApp); err != nil {
		return "", err
	}

	return refNumber, nil
}

func normalizeDateStr(s string) string {
	if len(s) >= 10 {
		return s[:10]
	}
	return s
}

func (u *applicationUsecase) GetStatus(citizenID string, birthDate string, refNumber string) (*entities.StatusResponseDTO, error) {
	app, err := u.repo.FindByCitizenID(citizenID)
	if err != nil {
		return nil, errors.New("APPLICATION_NOT_FOUND")
	}

	// Normalize: SQLite date อาจออกมาเป็น "1985-02-14T00:00:00Z" ต้องตัดให้เหลือ YYYY-MM-DD
	appBirthDate := normalizeDateStr(app.BirthDate)
	isMatchDOB := birthDate != "" && appBirthDate == normalizeDateStr(birthDate)
	isMatchRef := refNumber != "" && strings.EqualFold(app.ReferenceNumber, refNumber)

	if !isMatchDOB && !isMatchRef {
		return nil, errors.New("VERIFICATION_FAILED")
	}

	return &entities.StatusResponseDTO{
		ReferenceNumber: app.ReferenceNumber,
		CitizenID:       maskCitizenID(app.CitizenID),
		FullName:        maskName(app.FullName),
		Status:          app.Status,
		Reason:          app.Reason,
	}, nil
}

func (u *applicationUsecase) GetOfficerApplications(statusFilter string) ([]entities.Application, error) {
	return u.repo.FindAll(statusFilter)
}

func (u *applicationUsecase) ReviewApplication(id string, dto entities.UpdateStatusDTO) (*entities.Application, error) {
	app, err := u.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("APPLICATION_NOT_FOUND")
	}

	if dto.Status == "rejected" && strings.TrimSpace(dto.Reason) == "" {
		return nil, errors.New("REASON_REQUIRED_FOR_REJECT")
	}

	app.Status = dto.Status
	if dto.Status == "rejected" {
		app.Reason = strings.TrimSpace(dto.Reason)
	} else {
		app.Reason = ""
	}
	app.UpdatedAt = time.Now()

	if err := u.repo.Update(app); err != nil {
		return nil, err
	}

	return app, nil
}

package repositories

import (
	"errors"
	"strings"
	"sync"
	"time"
	"welfare-registration-backend/internal/entities"

	"github.com/google/uuid"
)

type inMemoryRepository struct {
	sync.RWMutex
	applications []entities.Application
}

func NewInMemoryRepository() entities.ApplicationRepository {
	return &inMemoryRepository{
		applications: []entities.Application{
			{
				ID:              "a1000000-0000-0000-0000-000000000001",
				ReferenceNumber: "WRS-2026-000001",
				CitizenID:       "1100400123450",
				FullName:        "สมชาย ใจดี",
				BirthDate:       "1990-05-12",
				AnnualIncome:    54000,
				CurrentAddress:  "99/1 แขวงดินแดง เขตดินแดง กทม. 10400",
				Status:          "pending",
				CreatedAt:       time.Now().Add(-48 * time.Hour),
				UpdatedAt:       time.Now().Add(-48 * time.Hour),
			},
			{
				ID:              "a2000000-0000-0000-0000-000000000002",
				ReferenceNumber: "WRS-2026-000002",
				CitizenID:       "1101700205673",
				FullName:        "สมหญิง รักดี",
				BirthDate:       "1988-11-20",
				AnnualIncome:    72000,
				CurrentAddress:  "12/4 ต.สุเทพ อ.เมือง จ.เชียงใหม่ 50200",
				Status:          "approved",
				CreatedAt:       time.Now().Add(-24 * time.Hour),
				UpdatedAt:       time.Now().Add(-24 * time.Hour),
			},
			{
				ID:              "a3000000-0000-0000-0000-000000000003",
				ReferenceNumber: "WRS-2026-000003",
				CitizenID:       "1102500337895",
				FullName:        "วิชัย ตั้งใจ",
				BirthDate:       "1985-02-14",
				AnnualIncome:    142000,
				CurrentAddress:  "45 ถ.มิตรภาพ ต.ในเมือง อ.เมือง จ.ขอนแก่น 40000",
				Status:          "rejected",
				Reason:          "รายได้รวมของครัวเรือนเกินเกณฑ์ที่กำหนด (ไม่เกิน 100,000 บาท/ปี)",
				CreatedAt:       time.Now().Add(-12 * time.Hour),
				UpdatedAt:       time.Now().Add(-12 * time.Hour),
			},
		},
	}
}

func (r *inMemoryRepository) FindAll(statusFilter string) ([]entities.Application, error) {
	r.RLock()
	defer r.RUnlock()

	var result []entities.Application
	for _, app := range r.applications {
		if statusFilter == "" || statusFilter == "all" || app.Status == statusFilter {
			result = append(result, app)
		}
	}
	if result == nil {
		result = []entities.Application{}
	}
	return result, nil
}

func (r *inMemoryRepository) FindByCitizenID(citizenID string) (*entities.Application, error) {
	r.RLock()
	defer r.RUnlock()

	cleanID := strings.ReplaceAll(citizenID, "-", "")
	for _, app := range r.applications {
		if app.CitizenID == cleanID {
			return &app, nil
		}
	}
	return nil, errors.New("APPLICATION_NOT_FOUND")
}

func (r *inMemoryRepository) FindByID(id string) (*entities.Application, error) {
	r.RLock()
	defer r.RUnlock()

	for _, app := range r.applications {
		if app.ID == id {
			return &app, nil
		}
	}
	return nil, errors.New("APPLICATION_NOT_FOUND")
}

func (r *inMemoryRepository) Save(app *entities.Application) error {
	r.Lock()
	defer r.Unlock()

	if app.ID == "" {
		app.ID = uuid.New().String()
	}
	r.applications = append([]entities.Application{*app}, r.applications...)
	return nil
}

func (r *inMemoryRepository) Update(app *entities.Application) error {
	r.Lock()
	defer r.Unlock()

	for i, existing := range r.applications {
		if existing.ID == app.ID {
			r.applications[i] = *app
			return nil
		}
	}
	return errors.New("APPLICATION_NOT_FOUND")
}

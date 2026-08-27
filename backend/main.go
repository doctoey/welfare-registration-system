package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

type Application struct {
	ID              int     `json:"id"`
	ReferenceNumber string  `json:"referenceNumber"`
	CitizenID       string  `json:"citizenId"`
	FullName        string  `json:"fullName"`
	BirthDate       string  `json:"birthDate"`
	AnnualIncome    float64 `json:"annualIncome"`
	CurrentAddress  string  `json:"currentAddress"`
	Status          string  `json:"status"`
	Reason          string  `json:"reason,omitempty"`
}

var applications = []Application{
	{
		ID:              1,
		ReferenceNumber: "WRS-2026-000001",
		CitizenID:       "1100400123450",
		FullName:        "สมชาย ใจดี",
		BirthDate:       "1990-05-12",
		AnnualIncome:    54000,
		CurrentAddress:  "99/1 แขวงดินแดง เขตดินแดง กทม. 10400",
		Status:          "pending",
	},
	{
		ID:              2,
		ReferenceNumber: "WRS-2026-000002",
		CitizenID:       "1101700205673",
		FullName:        "สมหญิง รักดี",
		BirthDate:       "1988-11-20",
		AnnualIncome:    72000,
		CurrentAddress:  "12/4 ต.สุเทพ อ.เมือง จ.เชียงใหม่ 50200",
		Status:          "approved",
	},
	{
		ID:              3,
		ReferenceNumber: "WRS-2026-000003",
		CitizenID:       "1102500337895",
		FullName:        "วิชัย ตั้งใจ",
		BirthDate:       "1985-02-14",
		AnnualIncome:    142000,
		CurrentAddress:  "45 ถ.มิตรภาพ ต.ในเมือง อ.เมือง จ.ขอนแก่น 40000",
		Status:          "rejected",
		Reason:          "รายได้รวมของครัวเรือนเกินเกณฑ์ที่กำหนด (ไม่เกิน 100,000 บาท/ปี)",
	},
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

func handleGetOfficer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	statusFilter := r.URL.Query().Get("status")
	var result []Application

	for _, app := range applications {
		if statusFilter == "" || statusFilter == "all" || app.Status == statusFilter {
			result = append(result, app)
		}
	}

	if result == nil {
		result = []Application{}
	}

	json.NewEncoder(w).Encode(result)
}

func handleGetStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	citizenID := r.PathValue("citizenId")

	for _, app := range applications {
		if app.CitizenID == citizenID {
			json.NewEncoder(w).Encode(app)
			return
		}
	}

	http.Error(w, "APPLICATION_NOT_FOUND", http.StatusNotFound)
}

func handleCreateApplication(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req Application
	json.NewDecoder(r.Body).Decode(&req)

	if req.CitizenID == "" || req.FullName == "" {
		http.Error(w, "กรุณากรอกข้อมูลให้ครบถ้วน", http.StatusBadRequest)
		return
	}

	if !isValidThaiID(req.CitizenID) {
		http.Error(w, "Invalid citizen ID checksum", http.StatusBadRequest)
		return
	}

	for _, app := range applications {
		if app.CitizenID == req.CitizenID {
			http.Error(w, "CITIZEN_ALREADY_REGISTERED", http.StatusConflict)
			return
		}
	}

	refNumber := fmt.Sprintf("WRS-2026-%06d", rand.Intn(900000)+100000)
	req.ID = len(applications) + 1
	req.ReferenceNumber = refNumber
	req.Status = "pending"

	applications = append([]Application{req}, applications...)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"referenceNumber": refNumber})
}

func handleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id, _ := strconv.Atoi(r.PathValue("id"))

	var body struct {
		Status string `json:"status"`
		Reason string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	if body.Status == "rejected" && strings.TrimSpace(body.Reason) == "" {
		http.Error(w, "REASON_REQUIRED_FOR_REJECT", http.StatusBadRequest)
		return
	}

	for i, app := range applications {
		if app.ID == id {
			applications[i].Status = body.Status
			applications[i].Reason = body.Reason

			json.NewEncoder(w).Encode(applications[i])
			return
		}
	}

	http.Error(w, "APPLICATION_NOT_FOUND", http.StatusNotFound)
}

func main() {
	http.HandleFunc("GET /api/v1/officer/applications", handleGetOfficer)
	http.HandleFunc("GET /api/v1/applications/status/{citizenId}", handleGetStatus)
	http.HandleFunc("POST /api/v1/applications", handleCreateApplication)
	http.HandleFunc("PATCH /api/v1/officer/applications/{id}/status", handleUpdateStatus)

	fmt.Println("port run at : http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

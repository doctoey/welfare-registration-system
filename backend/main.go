package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func handleGetOfficer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(applications)
}

func main() {
	http.HandleFunc("GET /api/v1/officer/applications", handleGetOfficer)

	fmt.Println("port run at : http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

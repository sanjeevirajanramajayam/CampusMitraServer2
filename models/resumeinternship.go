package models

type Internship struct {
	CompanyName string `json:"company_name"`
	Roll string `json:"roll"`
	Duration string `json:"duration"`
	IsPaid bool `json:"is_paid"`
	Domain string `json:"domain"`
}
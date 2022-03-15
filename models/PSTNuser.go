package models

type PstnUser struct {
	Email       string   `json:"email"`
	PhoneNumber string   `json:"phoneNumber"`
	SessionIds  []string `json:"sessionId"`
	HearingIds  []string `json:"hearingId"`
}

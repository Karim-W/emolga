package models

type PstnUser struct {
	Email       string   `json:"email"`
	PhoneNumber string   `json:"phoneNumber"`
	SessionIds  []string `json:"sessionId,omitempty"`
	HearingIds  []string `json:"hearingId,omitempty"`
}

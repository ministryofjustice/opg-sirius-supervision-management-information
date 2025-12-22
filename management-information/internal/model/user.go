package model

import (
	"slices"
)

type User struct {
	Id          int      `json:"id"`
	Name        string   `json:"displayName"`
	PhoneNumber string   `json:"phoneNumber"`
	Deleted     bool     `json:"deleted"`
	Email       string   `json:"email"`
	Firstname   string   `json:"firstname"`
	Surname     string   `json:"surname"`
	Roles       []string `json:"roles"`
	Locked      bool     `json:"locked"`
	Suspended   bool     `json:"suspended"`
}

func (m User) IsReportingUser() bool {
	return slices.Contains(m.Roles, "Reporting User")
}

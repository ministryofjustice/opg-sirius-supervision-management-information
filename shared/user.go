package shared

const (
	RoleReportingUser    = "Reporting User"
	RoleFinanceReporting = "Finance Reporting"
	RoleAny              = ""
)

type User struct {
	ID          int32    `json:"id"`
	DisplayName string   `json:"displayName"`
	Roles       []string `json:"roles"`
}

func (u User) IsReportingUser() bool {
	return contains(u.Roles, RoleReportingUser)
}

func (u User) HasRole(role string) bool {
	if role == RoleAny {
		return true
	}
	return contains(u.Roles, role)
}

func contains(haystack []string, needle string) bool {
	for _, h := range haystack {
		if h == needle {
			return true
		}
	}
	return false
}

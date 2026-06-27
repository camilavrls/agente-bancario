package policy

type Role string

const (
	RoleCustomer Role = "customer"
	RoleManager  Role = "manager"
	RoleAdmin    Role = "admin"
)

type AuthenticatedUser struct {
	ID         string
	Name       string
	Role       Role
	CustomerID string
}

type PolicyDecision struct {
	Allowed bool
	Reason  string
}

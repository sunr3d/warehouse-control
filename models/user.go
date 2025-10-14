package models

const (
	RoleAdmin   = "admin"
	RoleManager = "manager"
	RoleViewer  = "viewer"
)

type User struct {
	ID           int
	Username     string
	PasswordHash string
	Role         string
}

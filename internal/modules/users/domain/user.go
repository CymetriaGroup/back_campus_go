package domain

import "time"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type UserStatus string

const (
	StatusActive    UserStatus = "ACTIVE"
	StatusInactive  UserStatus = "INACTIVE"
	StatusSuspended UserStatus = "SUSPENDED"
	StatusBlocked   UserStatus = "BLOCKED"
)

type User struct {
	ID        string
	AuthID    string
	TenantID  string
	Name      string
	Email     string
	Role      Role
	Status    UserStatus
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

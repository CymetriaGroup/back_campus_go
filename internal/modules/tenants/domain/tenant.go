package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrTenantNotFound      = errors.New("tenant not found")
	ErrTenantAlreadyExists = errors.New("tenant already exists")
	ErrInvalidTenantData   = errors.New("invalid tenant data")
)

type TenantStatus string

const (
	StatusActive   TenantStatus = "ACTIVE"
	StatusInactive TenantStatus = "INACTIVE"
	StatusPending  TenantStatus = "PENDING"
)

type Tenant struct {
	ID        string
	Name      string
	Branding  map[string]any
	Domains   []string
	Limits    map[string]any
	Features  map[string]any
	Status    TenantStatus
	Settings  TenantSettings
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TenantSettings struct {
	TenantID           string
	Timezone           string
	Language           string
	DateFormat         string
	InstitutionalEmail string
	Policies           map[string]any
	FeatureFlags       map[string]any
}

type TenantRepository interface {
	Create(ctx context.Context, tenant *Tenant) error
	FindByID(ctx context.Context, id string) (*Tenant, error)
	FindByDomain(ctx context.Context, domain string) (*Tenant, error)
	List(ctx context.Context) ([]*Tenant, error)
	Update(ctx context.Context, tenant *Tenant) error
	Delete(ctx context.Context, id string) error
}

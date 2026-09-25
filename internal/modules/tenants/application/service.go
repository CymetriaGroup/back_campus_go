package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	"hexagonal-go-backend/internal/modules/tenants/domain"
)

type CreateTenantInput struct {
	ID                 string         `json:"id"`
	Name               string         `json:"name"`
	Branding           map[string]any `json:"branding"`
	Domains            []string       `json:"domains"`
	Limits             map[string]any `json:"limits"`
	Features           map[string]any `json:"features"`
	Timezone           string         `json:"timezone"`
	Language           string         `json:"language"`
	DateFormat         string         `json:"date_format"`
	InstitutionalEmail string         `json:"institutional_email"`
	Policies           map[string]any `json:"policies"`
	FeatureFlags       map[string]any `json:"feature_flags"`
}

type UpdateTenantInput struct {
	Name               *string        `json:"name"`
	Branding           map[string]any `json:"branding"`
	Domains            []string       `json:"domains"`
	Limits             map[string]any `json:"limits"`
	Features           map[string]any `json:"features"`
	Status             *string        `json:"status"`
	Timezone           *string        `json:"timezone"`
	Language           *string        `json:"language"`
	DateFormat         *string        `json:"date_format"`
	InstitutionalEmail *string        `json:"institutional_email"`
	Policies           map[string]any `json:"policies"`
	FeatureFlags       map[string]any `json:"feature_flags"`
}

type TenantService interface {
	Create(ctx context.Context, input CreateTenantInput) (*domain.Tenant, error)
	GetByID(ctx context.Context, id string) (*domain.Tenant, error)
	GetByDomain(ctx context.Context, domainName string) (*domain.Tenant, error)
	List(ctx context.Context) ([]*domain.Tenant, error)
	Update(ctx context.Context, id string, input UpdateTenantInput) (*domain.Tenant, error)
	Delete(ctx context.Context, id string) error
}

type tenantService struct {
	repo domain.TenantRepository
}

func NewTenantService(repo domain.TenantRepository) TenantService {
	return &tenantService{repo: repo}
}

func (s *tenantService) Create(ctx context.Context, input CreateTenantInput) (*domain.Tenant, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, fmt.Errorf("%w: name is required", domain.ErrInvalidTenantData)
	}

	tz := input.Timezone
	if tz == "" {
		tz = "UTC"
	}
	lang := input.Language
	if lang == "" {
		lang = "es"
	}
	df := input.DateFormat
	if df == "" {
		df = "YYYY-MM-DD"
	}

	now := time.Now()
	tenant := &domain.Tenant{
		ID:       input.ID,
		Name:     name,
		Branding: input.Branding,
		Domains:  input.Domains,
		Limits:   input.Limits,
		Features: input.Features,
		Status:   domain.StatusActive,
		Settings: domain.TenantSettings{
			TenantID:           input.ID,
			Timezone:           tz,
			Language:           lang,
			DateFormat:         df,
			InstitutionalEmail: input.InstitutionalEmail,
			Policies:           input.Policies,
			FeatureFlags:       input.FeatureFlags,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, tenant); err != nil {
		return nil, err
	}

	return tenant, nil
}

func (s *tenantService) GetByID(ctx context.Context, id string) (*domain.Tenant, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", domain.ErrInvalidTenantData)
	}
	return s.repo.FindByID(ctx, id)
}

func (s *tenantService) GetByDomain(ctx context.Context, domainName string) (*domain.Tenant, error) {
	domainName = strings.TrimSpace(domainName)
	if domainName == "" {
		return nil, fmt.Errorf("%w: domain is required", domain.ErrInvalidTenantData)
	}
	return s.repo.FindByDomain(ctx, domainName)
}

func (s *tenantService) List(ctx context.Context) ([]*domain.Tenant, error) {
	return s.repo.List(ctx)
}

func (s *tenantService) Update(ctx context.Context, id string, input UpdateTenantInput) (*domain.Tenant, error) {
	tenant, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, fmt.Errorf("%w: name cannot be empty", domain.ErrInvalidTenantData)
		}
		tenant.Name = name
	}
	if input.Branding != nil {
		tenant.Branding = input.Branding
	}
	if input.Domains != nil {
		tenant.Domains = input.Domains
	}
	if input.Limits != nil {
		tenant.Limits = input.Limits
	}
	if input.Features != nil {
		tenant.Features = input.Features
	}
	if input.Status != nil {
		tenant.Status = domain.TenantStatus(*input.Status)
	}
	if input.Timezone != nil {
		tenant.Settings.Timezone = *input.Timezone
	}
	if input.Language != nil {
		tenant.Settings.Language = *input.Language
	}
	if input.DateFormat != nil {
		tenant.Settings.DateFormat = *input.DateFormat
	}
	if input.InstitutionalEmail != nil {
		tenant.Settings.InstitutionalEmail = *input.InstitutionalEmail
	}
	if input.Policies != nil {
		tenant.Settings.Policies = input.Policies
	}
	if input.FeatureFlags != nil {
		tenant.Settings.FeatureFlags = input.FeatureFlags
	}

	tenant.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, tenant); err != nil {
		return nil, err
	}

	return tenant, nil
}

func (s *tenantService) Delete(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("%w: id is required", domain.ErrInvalidTenantData)
	}
	return s.repo.Delete(ctx, id)
}

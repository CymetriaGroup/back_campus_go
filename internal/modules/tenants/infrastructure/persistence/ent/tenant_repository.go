package ent

import (
	"context"
	"errors"
	"fmt"

	entclient "hexagonal-go-backend/internal/ent"
	enttenants "hexagonal-go-backend/internal/ent/tenants"
	enttenantsettings "hexagonal-go-backend/internal/ent/tenantsettings"
	"hexagonal-go-backend/internal/modules/tenants/domain"
)

type TenantRepository struct {
	client *entclient.Client
}

func NewTenantRepository(client *entclient.Client) *TenantRepository {
	return &TenantRepository{client: client}
}

func (r *TenantRepository) Create(ctx context.Context, t *domain.Tenant) error {
	builder := r.client.Tenants.Create().
		SetName(t.Name).
		SetBranding(t.Branding).
		SetDomains(t.Domains).
		SetLimits(t.Limits).
		SetFeatures(t.Features).
		SetStatus(enttenants.Status(t.Status))

	if t.ID != "" {
		builder.SetID(t.ID)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return translateError(err)
	}
	t.ID = created.ID

	// Create 1-to-1 TenantSettings
	_, err = r.client.TenantSettings.Create().
		SetTenantID(created.ID).
		SetTimezone(t.Settings.Timezone).
		SetLanguage(t.Settings.Language).
		SetDateFormat(t.Settings.DateFormat).
		SetInstitutionalEmail(t.Settings.InstitutionalEmail).
		SetPolicies(t.Settings.Policies).
		SetFeatureFlags(t.Settings.FeatureFlags).
		SetTenant(created).
		Save(ctx)

	return translateError(err)
}

func (r *TenantRepository) FindByID(ctx context.Context, id string) (*domain.Tenant, error) {
	entity, err := r.client.Tenants.Query().
		Where(enttenants.IDEQ(id)).
		WithSettings().
		Only(ctx)
	if err != nil {
		return nil, translateError(err)
	}
	return toDomain(entity), nil
}

func (r *TenantRepository) FindByDomain(ctx context.Context, domainName string) (*domain.Tenant, error) {
	// Search in JSON array of domains
	tenantsList, err := r.client.Tenants.Query().
		WithSettings().
		All(ctx)
	if err != nil {
		return nil, translateError(err)
	}

	for _, t := range tenantsList {
		for _, d := range t.Domains {
			if d == domainName {
				return toDomain(t), nil
			}
		}
	}
	return nil, domain.ErrTenantNotFound
}

func (r *TenantRepository) List(ctx context.Context) ([]*domain.Tenant, error) {
	entities, err := r.client.Tenants.Query().
		WithSettings().
		All(ctx)
	if err != nil {
		return nil, translateError(err)
	}

	result := make([]*domain.Tenant, 0, len(entities))
	for _, e := range entities {
		result = append(result, toDomain(e))
	}
	return result, nil
}

func (r *TenantRepository) Update(ctx context.Context, t *domain.Tenant) error {
	_, err := r.client.Tenants.UpdateOneID(t.ID).
		SetName(t.Name).
		SetBranding(t.Branding).
		SetDomains(t.Domains).
		SetLimits(t.Limits).
		SetFeatures(t.Features).
		SetStatus(enttenants.Status(t.Status)).
		Save(ctx)
	if err != nil {
		return translateError(err)
	}

	if t.Settings.TenantID != "" {
		settingsEntity, err := r.client.TenantSettings.Query().
			Where(enttenantsettings.HasTenantWith(enttenants.IDEQ(t.ID))).
			Only(ctx)
		if err == nil && settingsEntity != nil {
			_, _ = r.client.TenantSettings.UpdateOne(settingsEntity).
				SetTimezone(t.Settings.Timezone).
				SetLanguage(t.Settings.Language).
				SetDateFormat(t.Settings.DateFormat).
				SetInstitutionalEmail(t.Settings.InstitutionalEmail).
				SetPolicies(t.Settings.Policies).
				SetFeatureFlags(t.Settings.FeatureFlags).
				Save(ctx)
		}
	}

	return nil
}

func (r *TenantRepository) Delete(ctx context.Context, id string) error {
	return translateError(r.client.Tenants.DeleteOneID(id).Exec(ctx))
}

func toDomain(e *entclient.Tenants) *domain.Tenant {
	if e == nil {
		return nil
	}
	t := &domain.Tenant{
		ID:        e.ID,
		Name:      e.Name,
		Branding:  e.Branding,
		Domains:   e.Domains,
		Limits:    e.Limits,
		Features:  e.Features,
		Status:    domain.TenantStatus(e.Status),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}

	if e.Edges.Settings != nil {
		t.Settings = domain.TenantSettings{
			TenantID:           e.ID,
			Timezone:           e.Edges.Settings.Timezone,
			Language:           e.Edges.Settings.Language,
			DateFormat:         e.Edges.Settings.DateFormat,
			InstitutionalEmail: e.Edges.Settings.InstitutionalEmail,
			Policies:           e.Edges.Settings.Policies,
			FeatureFlags:       e.Edges.Settings.FeatureFlags,
		}
	}
	return t
}

func translateError(err error) error {
	if err == nil {
		return nil
	}
	if entclient.IsNotFound(err) {
		return domain.ErrTenantNotFound
	}
	if entclient.IsConstraintError(err) {
		return domain.ErrTenantAlreadyExists
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	return fmt.Errorf("tenant repository: %w", err)
}

var _ domain.TenantRepository = (*TenantRepository)(nil)


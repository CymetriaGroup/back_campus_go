package postgres

import (
	"context"
	"errors"
	"fmt"

	entclient "hexagonal-go-backend/internal/ent"
	"hexagonal-go-backend/internal/ent/coursecategories"
	"hexagonal-go-backend/internal/ent/courseversions"
	"hexagonal-go-backend/internal/ent/syllabi"
	"hexagonal-go-backend/internal/modules/courses/application"
	"hexagonal-go-backend/internal/modules/courses/domain"
	"hexagonal-go-backend/internal/modules/courses/infrastructure/persistence/postgres/mappers"
	"hexagonal-go-backend/internal/platform/identifier"
)

type Repository struct{ client *entclient.Client }

func NewRepository(client *entclient.Client) *Repository { return &Repository{client: client} }

var _ application.Repository = (*Repository)(nil)

func translate(err error) error {
	if err == nil {
		return nil
	}
	if entclient.IsNotFound(err) {
		return fmt.Errorf("%w: %v", domain.ErrNotFound, err)
	}
	if entclient.IsConstraintError(err) {
		return fmt.Errorf("%w: %v", domain.ErrConflict, err)
	}
	if entclient.IsValidationError(err) {
		return fmt.Errorf("%w: %v", domain.ErrInvalid, err)
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	return fmt.Errorf("courses persistence: %w", err)
}

func (r *Repository) CreateCategory(ctx context.Context, c domain.Category) (domain.Category, error) {
	if err := c.Validate(); err != nil {
		return domain.Category{}, err
	}
	if c.ParentID != "" {
		parent, err := r.client.CourseCategories.Query().Where(coursecategories.IDEQ(c.ParentID), coursecategories.TenantIDEQ(c.TenantID)).Only(ctx)
		if err != nil {
			return domain.Category{}, translate(err)
		}
		_ = parent
	}
	id := c.ID
	if id == "" {
		id = identifier.New()
	}
	parentID := c.ParentID
	// The existing schema requires a non-null ULID parent_id; a root points to itself.
	if parentID == "" {
		parentID = id
	}
	row, err := r.client.CourseCategories.Create().SetID(id).SetTenantID(c.TenantID).SetCode(c.Code).SetName(c.Name).SetDescription(c.Description).SetParentID(parentID).Save(ctx)
	if err != nil {
		return domain.Category{}, translate(err)
	}
	return mappers.Category(row), nil
}

func (r *Repository) ListCategories(ctx context.Context, tenantID string) ([]domain.Category, error) {
	if tenantID == "" {
		return nil, domain.ErrInvalid
	}
	rows, err := r.client.CourseCategories.Query().Where(coursecategories.TenantIDEQ(tenantID)).Order(coursecategories.ByCode()).All(ctx)
	if err != nil {
		return nil, translate(err)
	}
	result := make([]domain.Category, 0, len(rows))
	for _, row := range rows {
		result = append(result, mappers.Category(row))
	}
	return result, nil
}

func (r *Repository) CreateTemplate(ctx context.Context, t domain.Template) (domain.Template, error) {
	if err := t.Validate(); err != nil {
		return domain.Template{}, err
	}
	create := r.client.CourseTemplates.Create().SetCode(t.Code).SetTitle(t.Title).SetDescription(t.Description)
	if t.ID != "" {
		create.SetID(t.ID)
	}
	row, err := create.Save(ctx)
	if err != nil {
		return domain.Template{}, translate(err)
	}
	return mappers.Template(row), nil
}

func (r *Repository) ListTemplates(ctx context.Context) ([]domain.Template, error) {
	rows, err := r.client.CourseTemplates.Query().All(ctx)
	if err != nil {
		return nil, translate(err)
	}
	result := make([]domain.Template, 0, len(rows))
	for _, row := range rows {
		result = append(result, mappers.Template(row))
	}
	return result, nil
}

func (r *Repository) GetTemplate(ctx context.Context, id string) (domain.Template, error) {
	if id == "" {
		return domain.Template{}, domain.ErrInvalid
	}
	row, err := r.client.CourseTemplates.Get(ctx, id)
	if err != nil {
		return domain.Template{}, translate(err)
	}
	return mappers.Template(row), nil
}

func (r *Repository) CreateVersion(ctx context.Context, v domain.Version) (domain.Version, error) {
	if err := v.Validate(); err != nil {
		return domain.Version{}, err
	}
	_, err := r.client.CourseTemplates.Get(ctx, v.TemplateID)
	if err != nil {
		return domain.Version{}, translate(err)
	}
	// Ent does not declare a unique (template_id, version_tag) index.
	exists, err := r.client.CourseVersions.Query().Where(courseversions.TemplateIDEQ(v.TemplateID), courseversions.VersionTagEQ(v.Tag)).Exist(ctx)
	if err != nil {
		return domain.Version{}, translate(err)
	}
	if exists {
		return domain.Version{}, domain.ErrConflict
	}
	create := r.client.CourseVersions.Create().SetTemplateID(v.TemplateID).SetVersionTag(v.Tag).SetStatus(v.Status).SetEstimatedHours(v.EstimatedHours)
	if v.ID != "" {
		create.SetID(v.ID)
	}
	row, err := create.Save(ctx)
	if err != nil {
		return domain.Version{}, translate(err)
	}
	return mappers.Version(row), nil
}

func (r *Repository) ListVersions(ctx context.Context, templateID string) ([]domain.Version, error) {
	if templateID == "" {
		return nil, domain.ErrInvalid
	}
	_, err := r.client.CourseTemplates.Get(ctx, templateID)
	if err != nil {
		return nil, translate(err)
	}
	rows, err := r.client.CourseVersions.Query().Where(courseversions.TemplateIDEQ(templateID)).Order(courseversions.ByVersionTag()).All(ctx)
	if err != nil {
		return nil, translate(err)
	}
	result := make([]domain.Version, 0, len(rows))
	for _, row := range rows {
		result = append(result, mappers.Version(row))
	}
	return result, nil
}

func (r *Repository) SetSyllabus(ctx context.Context, s domain.Syllabus) (domain.Syllabus, error) {
	if err := s.Validate(); err != nil {
		return domain.Syllabus{}, err
	}
	_, err := r.client.CourseVersions.Get(ctx, s.VersionID)
	if err != nil {
		return domain.Syllabus{}, translate(err)
	}
	rows, err := r.client.Syllabi.Query().Where(syllabi.VersionIDEQ(s.VersionID)).Limit(2).All(ctx)
	if err != nil {
		return domain.Syllabus{}, translate(err)
	}
	if len(rows) > 1 {
		return domain.Syllabus{}, domain.ErrConflict
	}
	if len(rows) == 1 {
		if s.ID != "" && s.ID != rows[0].ID {
			return domain.Syllabus{}, domain.ErrConflict
		}
		row, err := r.client.Syllabi.UpdateOneID(rows[0].ID).SetObjectives(s.Objectives).SetEntryProfile(s.EntryProfile).SetExitProfile(s.ExitProfile).SetMethodology(s.Methodology).SetDurationsHours(s.DurationHours).Save(ctx)
		if err != nil {
			return domain.Syllabus{}, translate(err)
		}
		return mappers.Syllabus(row), nil
	}
	create := r.client.Syllabi.Create().SetVersionID(s.VersionID).SetObjectives(s.Objectives).SetEntryProfile(s.EntryProfile).SetExitProfile(s.ExitProfile).SetMethodology(s.Methodology).SetDurationsHours(s.DurationHours)
	if s.ID != "" {
		create.SetID(s.ID)
	}
	row, err := create.Save(ctx)
	if err != nil {
		return domain.Syllabus{}, translate(err)
	}
	return mappers.Syllabus(row), nil
}

func (r *Repository) GetSyllabus(ctx context.Context, versionID string) (domain.Syllabus, error) {
	if versionID == "" {
		return domain.Syllabus{}, domain.ErrInvalid
	}
	rows, err := r.client.Syllabi.Query().Where(syllabi.VersionIDEQ(versionID)).Limit(2).All(ctx)
	if err != nil {
		return domain.Syllabus{}, translate(err)
	}
	if len(rows) > 1 {
		return domain.Syllabus{}, domain.ErrConflict
	}
	if len(rows) == 0 {
		return domain.Syllabus{}, domain.ErrNotFound
	}
	return mappers.Syllabus(rows[0]), nil
}

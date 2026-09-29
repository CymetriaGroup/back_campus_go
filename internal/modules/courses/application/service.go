package application

import (
	"context"
	"hexagonal-go-backend/internal/modules/courses/domain"
)

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func (s *Service) CreateCourse(ctx context.Context, item domain.Course) (domain.Course, error) {
	if err := item.Validate(); err != nil {
		return domain.Course{}, err
	}
	return s.repository.CreateCourse(ctx, item)
}
func (s *Service) GetCourse(ctx context.Context, versionID string) (domain.Course, error) {
	if versionID == "" {
		return domain.Course{}, domain.ErrInvalid
	}
	return s.repository.GetCourse(ctx, versionID)
}
func (s *Service) CreateCategory(ctx context.Context, item domain.Category) (domain.Category, error) {
	if err := item.Validate(); err != nil {
		return domain.Category{}, err
	}
	return s.repository.CreateCategory(ctx, item)
}
func (s *Service) ListCategories(ctx context.Context, tenantID string) ([]domain.Category, error) {
	return s.repository.ListCategories(ctx, tenantID)
}
func (s *Service) CreateTemplate(ctx context.Context, item domain.Template) (domain.Template, error) {
	if err := item.Validate(); err != nil {
		return domain.Template{}, err
	}
	return s.repository.CreateTemplate(ctx, item)
}
func (s *Service) ListTemplates(ctx context.Context) ([]domain.Template, error) {
	return s.repository.ListTemplates(ctx)
}
func (s *Service) GetTemplate(ctx context.Context, id string) (domain.Template, error) {
	return s.repository.GetTemplate(ctx, id)
}
func (s *Service) CreateVersion(ctx context.Context, item domain.Version) (domain.Version, error) {
	if err := item.Validate(); err != nil {
		return domain.Version{}, err
	}
	return s.repository.CreateVersion(ctx, item)
}
func (s *Service) ListVersions(ctx context.Context, id string) ([]domain.Version, error) {
	return s.repository.ListVersions(ctx, id)
}
func (s *Service) SetSyllabus(ctx context.Context, item domain.Syllabus) (domain.Syllabus, error) {
	if err := item.Validate(); err != nil {
		return domain.Syllabus{}, err
	}
	return s.repository.SetSyllabus(ctx, item)
}
func (s *Service) GetSyllabus(ctx context.Context, id string) (domain.Syllabus, error) {
	return s.repository.GetSyllabus(ctx, id)
}

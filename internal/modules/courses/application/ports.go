package application

import (
	"context"
	"hexagonal-go-backend/internal/modules/courses/domain"
)

type Repository interface {
	CreateCourse(context.Context, domain.Course) (domain.Course, error)
	GetCourse(context.Context, string) (domain.Course, error)
	CreateCategory(context.Context, domain.Category) (domain.Category, error)
	ListCategories(context.Context, string) ([]domain.Category, error)
	CreateTemplate(context.Context, domain.Template) (domain.Template, error)
	ListTemplates(context.Context) ([]domain.Template, error)
	GetTemplate(context.Context, string) (domain.Template, error)
	CreateVersion(context.Context, domain.Version) (domain.Version, error)
	ListVersions(context.Context, string) ([]domain.Version, error)
	SetSyllabus(context.Context, domain.Syllabus) (domain.Syllabus, error)
	GetSyllabus(context.Context, string) (domain.Syllabus, error)
}

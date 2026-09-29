package v1

import "hexagonal-go-backend/internal/modules/courses/domain"

type CreateCourseRequest struct {
	Template CreateTemplateRequest `json:"template" binding:"required"`
	Version  CreateVersionRequest  `json:"version" binding:"required"`
	Syllabus *SetSyllabusRequest   `json:"syllabus"`
	Modules  []CreateModuleRequest `json:"modules" binding:"dive"`
}

type CreateModuleRequest struct {
	Title         string                `json:"title" binding:"required"`
	SequenceOrder int                   `json:"sequence_order"`
	Lessons       []CreateLessonRequest `json:"lessons" binding:"dive"`
}

type CreateLessonRequest struct {
	Title         string                  `json:"title" binding:"required"`
	SequenceOrder int                     `json:"sequence_order"`
	Activities    []CreateActivityRequest `json:"activities" binding:"dive"`
}

type CreateActivityRequest struct {
	Title         string                  `json:"title" binding:"required"`
	Type          string                  `json:"type" binding:"required"`
	IsRequired    *bool                   `json:"is_required"`
	SequenceOrder int                     `json:"sequence_order"`
	Resources     []CreateResourceRequest `json:"resources" binding:"dive"`
}

type CreateResourceRequest struct {
	Type          string         `json:"type" binding:"required"`
	Name          string         `json:"name" binding:"required"`
	URLStorageKey string         `json:"url_storage_key" binding:"required"`
	MIMEType      string         `json:"mime_type"`
	SizeBytes     int            `json:"size_bytes" binding:"min=0"`
	Position      int            `json:"position"`
	Metadata      map[string]any `json:"metadata"`
}

type CourseResponse struct {
	Template TemplateResponse  `json:"template"`
	Version  VersionResponse   `json:"version"`
	Syllabus *SyllabusResponse `json:"syllabus"`
	Modules  []ModuleResponse  `json:"modules"`
}

type ModuleResponse struct {
	ID            string           `json:"id"`
	Title         string           `json:"title"`
	SequenceOrder int              `json:"sequence_order"`
	Lessons       []LessonResponse `json:"lessons"`
}

type LessonResponse struct {
	ID            string             `json:"id"`
	Title         string             `json:"title"`
	SequenceOrder int                `json:"sequence_order"`
	Activities    []ActivityResponse `json:"activities"`
}

type ActivityResponse struct {
	ID            string             `json:"id"`
	Title         string             `json:"title"`
	Type          string             `json:"type"`
	IsRequired    bool               `json:"is_required"`
	SequenceOrder int                `json:"sequence_order"`
	Resources     []ResourceResponse `json:"resources"`
}

type ResourceResponse struct {
	ID            string         `json:"id"`
	Type          string         `json:"type"`
	Name          string         `json:"name"`
	URLStorageKey string         `json:"url_storage_key"`
	MIMEType      string         `json:"mime_type"`
	SizeBytes     int            `json:"size_bytes"`
	Position      int            `json:"position"`
	Metadata      map[string]any `json:"metadata"`
}

func (r CreateCourseRequest) toDomain() domain.Course {
	course := domain.Course{
		Template: domain.Template{Code: r.Template.Code, Title: r.Template.Title, Description: r.Template.Description},
		Version:  domain.Version{Tag: r.Version.Tag, Status: r.Version.Status, EstimatedHours: r.Version.EstimatedHours},
		Modules:  make([]domain.Module, len(r.Modules)),
	}
	if r.Syllabus != nil {
		course.Syllabus = &domain.Syllabus{Objectives: r.Syllabus.Objectives, EntryProfile: r.Syllabus.EntryProfile, ExitProfile: r.Syllabus.ExitProfile, Methodology: r.Syllabus.Methodology, DurationHours: r.Syllabus.DurationHours}
	}
	for i, m := range r.Modules {
		module := domain.Module{Title: m.Title, SequenceOrder: m.SequenceOrder, Lessons: make([]domain.Lesson, len(m.Lessons))}
		for j, l := range m.Lessons {
			lesson := domain.Lesson{Title: l.Title, SequenceOrder: l.SequenceOrder, Activities: make([]domain.Activity, len(l.Activities))}
			for k, a := range l.Activities {
				required := true
				if a.IsRequired != nil {
					required = *a.IsRequired
				}
				activity := domain.Activity{Title: a.Title, Type: a.Type, IsRequired: required, SequenceOrder: a.SequenceOrder, Resources: make([]domain.Resource, len(a.Resources))}
				for n, resource := range a.Resources {
					activity.Resources[n] = domain.Resource{Type: resource.Type, Name: resource.Name, URLStorageKey: resource.URLStorageKey, MIMEType: resource.MIMEType, SizeBytes: resource.SizeBytes, Position: resource.Position, Metadata: resource.Metadata}
				}
				lesson.Activities[k] = activity
			}
			module.Lessons[j] = lesson
		}
		course.Modules[i] = module
	}
	return course
}

func courseResponse(v domain.Course) CourseResponse {
	result := CourseResponse{Template: templateResponse(v.Template), Version: versionResponse(v.Version), Modules: make([]ModuleResponse, len(v.Modules))}
	if v.Syllabus != nil {
		syllabus := syllabusResponse(*v.Syllabus)
		result.Syllabus = &syllabus
	}
	for i, m := range v.Modules {
		module := ModuleResponse{ID: m.ID, Title: m.Title, SequenceOrder: m.SequenceOrder, Lessons: make([]LessonResponse, len(m.Lessons))}
		for j, l := range m.Lessons {
			lesson := LessonResponse{ID: l.ID, Title: l.Title, SequenceOrder: l.SequenceOrder, Activities: make([]ActivityResponse, len(l.Activities))}
			for k, a := range l.Activities {
				activity := ActivityResponse{ID: a.ID, Title: a.Title, Type: a.Type, IsRequired: a.IsRequired, SequenceOrder: a.SequenceOrder, Resources: make([]ResourceResponse, len(a.Resources))}
				for n, resource := range a.Resources {
					activity.Resources[n] = ResourceResponse{ID: resource.ID, Type: resource.Type, Name: resource.Name, URLStorageKey: resource.URLStorageKey, MIMEType: resource.MIMEType, SizeBytes: resource.SizeBytes, Position: resource.Position, Metadata: resource.Metadata}
				}
				lesson.Activities[k] = activity
			}
			module.Lessons[j] = lesson
		}
		result.Modules[i] = module
	}
	return result
}

type CreateCategoryRequest struct {
	TenantID    string `json:"tenant_id" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ParentID    string `json:"parent_id"`
}

type CategoryResponse struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenant_id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    string `json:"parent_id"`
}

type CreateTemplateRequest struct {
	Code        string `json:"code" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type TemplateResponse struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateVersionRequest struct {
	Tag            string `json:"tag" binding:"required"`
	Status         string `json:"status" binding:"required,oneof=DRAFT PUBLISHED ARCHIVED"`
	EstimatedHours int    `json:"estimated_hours" binding:"min=0"`
}

type VersionResponse struct {
	ID             string `json:"id"`
	TemplateID     string `json:"template_id"`
	Tag            string `json:"tag"`
	Status         string `json:"status"`
	EstimatedHours int    `json:"estimated_hours"`
}

type SetSyllabusRequest struct {
	Objectives    string `json:"objectives"`
	EntryProfile  string `json:"entry_profile"`
	ExitProfile   string `json:"exit_profile"`
	Methodology   string `json:"methodology"`
	DurationHours int    `json:"duration_hours" binding:"min=0"`
}

type SyllabusResponse struct {
	ID            string `json:"id"`
	VersionID     string `json:"version_id"`
	Objectives    string `json:"objectives"`
	EntryProfile  string `json:"entry_profile"`
	ExitProfile   string `json:"exit_profile"`
	Methodology   string `json:"methodology"`
	DurationHours int    `json:"duration_hours"`
}

func categoryResponse(v domain.Category) CategoryResponse {
	return CategoryResponse{ID: v.ID, TenantID: v.TenantID, Code: v.Code, Name: v.Name, Description: v.Description, ParentID: v.ParentID}
}
func templateResponse(v domain.Template) TemplateResponse {
	return TemplateResponse{v.ID, v.Code, v.Title, v.Description}
}
func versionResponse(v domain.Version) VersionResponse {
	return VersionResponse{v.ID, v.TemplateID, v.Tag, v.Status, v.EstimatedHours}
}
func syllabusResponse(v domain.Syllabus) SyllabusResponse {
	return SyllabusResponse{v.ID, v.VersionID, v.Objectives, v.EntryProfile, v.ExitProfile, v.Methodology, v.DurationHours}
}

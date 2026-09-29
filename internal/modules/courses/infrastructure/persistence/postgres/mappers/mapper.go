package mappers

import (
	"encoding/json"

	"hexagonal-go-backend/internal/ent"
	"hexagonal-go-backend/internal/modules/courses/domain"
)

func Category(row *ent.CourseCategories) domain.Category {
	parentID := row.ParentID
	if parentID == row.ID {
		parentID = ""
	}
	return domain.Category{ID: row.ID, TenantID: row.TenantID, Code: row.Code, Name: row.Name, Description: row.Description, ParentID: parentID}
}

func Template(row *ent.CourseTemplates) domain.Template {
	return domain.Template{ID: row.ID, Code: row.Code, Title: row.Title, Description: row.Description}
}

func Version(row *ent.CourseVersions) domain.Version {
	return domain.Version{ID: row.ID, TemplateID: row.TemplateID, Tag: row.VersionTag, Status: row.Status, EstimatedHours: row.EstimatedHours}
}

func Module(row *ent.CourseModules) domain.Module {
	return domain.Module{ID: row.ID, Title: row.Title, SequenceOrder: row.SequenceOrder}
}

func Lesson(row *ent.Lessons) domain.Lesson {
	return domain.Lesson{ID: row.ID, Title: row.Title, SequenceOrder: row.SequenceOrder}
}

func Activity(row *ent.Activities) domain.Activity {
	return domain.Activity{ID: row.ID, Title: row.Title, Type: row.Type, IsRequired: row.IsRequired, SequenceOrder: row.SequenceOrder}
}

func Resource(row *ent.CourseResources) (domain.Resource, error) {
	var metadata map[string]any
	if len(row.Metadata) > 0 {
		if err := json.Unmarshal(row.Metadata, &metadata); err != nil {
			return domain.Resource{}, err
		}
	}
	return domain.Resource{ID: row.ID, Type: row.Type, Name: row.Name, URLStorageKey: row.URLStorageKey, MIMEType: row.MimeType, SizeBytes: int(row.SizeBytes), Position: row.Position, Metadata: metadata}, nil
}

func Syllabus(row *ent.Syllabi) domain.Syllabus {
	return domain.Syllabus{ID: row.ID, VersionID: row.VersionID, Objectives: row.Objectives, EntryProfile: row.EntryProfile, ExitProfile: row.ExitProfile, Methodology: row.Methodology, DurationHours: row.DurationsHours}
}

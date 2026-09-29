package postgres

import (
	"context"
	"encoding/json"

	entclient "hexagonal-go-backend/internal/ent"
	"hexagonal-go-backend/internal/ent/activities"
	"hexagonal-go-backend/internal/ent/coursemodules"
	"hexagonal-go-backend/internal/ent/courseresources"
	"hexagonal-go-backend/internal/ent/courseversions"
	"hexagonal-go-backend/internal/ent/lessons"

	"hexagonal-go-backend/internal/modules/courses/domain"
	"hexagonal-go-backend/internal/modules/courses/infrastructure/persistence/postgres/mappers"
)

func (r *Repository) CreateCourse(ctx context.Context, c domain.Course) (result domain.Course, err error) {
	if err = c.Validate(); err != nil {
		return domain.Course{}, err
	}
	// The aggregate is new: supplied foreign keys must never redirect it to existing records.
	if c.Version.TemplateID != "" && c.Version.TemplateID != c.Template.ID {
		return domain.Course{}, domain.ErrInvalid
	}
	if c.Syllabus != nil && c.Syllabus.VersionID != "" && c.Syllabus.VersionID != c.Version.ID {
		return domain.Course{}, domain.ErrInvalid
	}
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return domain.Course{}, translate(err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			result = domain.Course{}
		}
	}()
	t := tx.CourseTemplates.Create().SetCode(c.Template.Code).SetTitle(c.Template.Title).SetDescription(c.Template.Description)
	if c.Template.ID != "" {
		t.SetID(c.Template.ID)
	}
	template, e := t.Save(ctx)
	if e != nil {
		err = translate(e)
		return
	}
	result.Template = mappers.Template(template)
	v := tx.CourseVersions.Create().SetTemplateID(template.ID).SetVersionTag(c.Version.Tag).SetStatus(c.Version.Status).SetEstimatedHours(c.Version.EstimatedHours)
	if c.Version.ID != "" {
		v.SetID(c.Version.ID)
	}
	version, e := v.Save(ctx)
	if e != nil {
		err = translate(e)
		return
	}
	result.Version = mappers.Version(version)
	if c.Syllabus != nil {
		s := tx.Syllabi.Create().SetVersionID(version.ID).SetObjectives(c.Syllabus.Objectives).SetEntryProfile(c.Syllabus.EntryProfile).SetExitProfile(c.Syllabus.ExitProfile).SetMethodology(c.Syllabus.Methodology).SetDurationsHours(c.Syllabus.DurationHours)
		if c.Syllabus.ID != "" {
			s.SetID(c.Syllabus.ID)
		}
		row, saveErr := s.Save(ctx)
		if saveErr != nil {
			err = translate(saveErr)
			return
		}
		mapped := mappers.Syllabus(row)
		result.Syllabus = &mapped
	}
	result.Modules = make([]domain.Module, 0, len(c.Modules))
	for _, m := range c.Modules {
		create := tx.CourseModules.Create().SetVersionID(version.ID).SetTitle(m.Title).SetSequenceOrder(m.SequenceOrder)
		if m.ID != "" {
			create.SetID(m.ID)
		}
		row, saveErr := create.Save(ctx)
		if saveErr != nil {
			err = translate(saveErr)
			return
		}
		mapped := mappers.Module(row)
		mapped.Lessons = make([]domain.Lesson, 0, len(m.Lessons))
		for _, l := range m.Lessons {
			create := tx.Lessons.Create().SetModuleID(row.ID).SetTitle(l.Title).SetSequenceOrder(l.SequenceOrder)
			if l.ID != "" {
				create.SetID(l.ID)
			}
			lesson, saveErr := create.Save(ctx)
			if saveErr != nil {
				err = translate(saveErr)
				return
			}
			mappedLesson := mappers.Lesson(lesson)
			mappedLesson.Activities = make([]domain.Activity, 0, len(l.Activities))
			for _, a := range l.Activities {
				create := tx.Activities.Create().SetLessonID(lesson.ID).SetTitle(a.Title).SetType(a.Type).SetIsRequired(a.IsRequired).SetSequenceOrder(a.SequenceOrder)
				if a.ID != "" {
					create.SetID(a.ID)
				}
				activity, saveErr := create.Save(ctx)
				if saveErr != nil {
					err = translate(saveErr)
					return
				}
				mappedActivity := mappers.Activity(activity)
				mappedActivity.Resources = make([]domain.Resource, 0, len(a.Resources))
				for _, resource := range a.Resources {
					metadata, marshalErr := json.Marshal(resource.Metadata)
					if marshalErr != nil {
						err = domain.ErrInvalid
						return
					}
					create := tx.CourseResources.Create().SetActivityID(activity.ID).SetType(resource.Type).SetName(resource.Name).SetURLStorageKey(resource.URLStorageKey).SetMimeType(resource.MIMEType).SetSizeBytes(int64(resource.SizeBytes)).SetPosition(resource.Position)
					if resource.Metadata != nil {
						create.SetMetadata(json.RawMessage(metadata))
					}
					if resource.ID != "" {
						create.SetID(resource.ID)
					}
					saved, saveErr := create.Save(ctx)
					if saveErr != nil {
						err = translate(saveErr)
						return
					}
					mappedResource, mapErr := mappers.Resource(saved)
					if mapErr != nil {
						err = translate(mapErr)
						return
					}
					mappedActivity.Resources = append(mappedActivity.Resources, mappedResource)
				}
				mappedLesson.Activities = append(mappedLesson.Activities, mappedActivity)
			}
			mapped.Lessons = append(mapped.Lessons, mappedLesson)
		}
		result.Modules = append(result.Modules, mapped)
	}
	if e = tx.Commit(); e != nil {
		err = translate(e)
		return
	}
	return result, nil
}

func (r *Repository) GetCourse(ctx context.Context, id string) (domain.Course, error) {
	if id == "" {
		return domain.Course{}, domain.ErrInvalid
	}
	version, err := r.client.CourseVersions.Query().Where(courseversions.IDEQ(id)).WithTemplate().WithSyllabi().WithModules(func(q *entclient.CourseModulesQuery) {
		q.Order(coursemodules.BySequenceOrder()).WithLessons(func(q *entclient.LessonsQuery) {
			q.Order(lessons.BySequenceOrder()).WithActivities(func(q *entclient.ActivitiesQuery) {
				q.Order(activities.BySequenceOrder()).WithResources(func(q *entclient.CourseResourcesQuery) { q.Order(courseresources.ByPosition()) })
			})
		})
	}).Only(ctx)
	if err != nil {
		return domain.Course{}, translate(err)
	}
	result := domain.Course{Template: mappers.Template(version.Edges.Template), Version: mappers.Version(version), Modules: make([]domain.Module, 0, len(version.Edges.Modules))}
	if len(version.Edges.Syllabi) > 1 {
		return domain.Course{}, domain.ErrConflict
	}
	if len(version.Edges.Syllabi) == 1 {
		s := mappers.Syllabus(version.Edges.Syllabi[0])
		result.Syllabus = &s
	}
	for _, m := range version.Edges.Modules {
		module := mappers.Module(m)
		module.Lessons = make([]domain.Lesson, 0, len(m.Edges.Lessons))
		for _, l := range m.Edges.Lessons {
			lesson := mappers.Lesson(l)
			lesson.Activities = make([]domain.Activity, 0, len(l.Edges.Activities))
			for _, a := range l.Edges.Activities {
				activity := mappers.Activity(a)
				activity.Resources = make([]domain.Resource, 0, len(a.Edges.Resources))
				for _, row := range a.Edges.Resources {
					resource, mapErr := mappers.Resource(row)
					if mapErr != nil {
						return domain.Course{}, translate(mapErr)
					}
					activity.Resources = append(activity.Resources, resource)
				}
				lesson.Activities = append(lesson.Activities, activity)
			}
			module.Lessons = append(module.Lessons, lesson)
		}
		result.Modules = append(result.Modules, module)
	}
	return result, nil
}

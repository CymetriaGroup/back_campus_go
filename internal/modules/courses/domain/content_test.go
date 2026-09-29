package domain

import (
	"errors"
	"testing"
)

func TestCourseValidate(t *testing.T) {
	valid := func() Course {
		return Course{
			Template: Template{Code: "GO", Title: "Go"}, Version: Version{Tag: "v1", Status: "DRAFT"},
			Modules: []Module{{Title: "Module", SequenceOrder: 1, Lessons: []Lesson{{Title: "Lesson", SequenceOrder: 1, Activities: []Activity{{Title: "Exercise", Type: "TASK", SequenceOrder: 1, IsRequired: true, Resources: []Resource{{Type: "FILE", Name: "Guide", URLStorageKey: "courses/guide.pdf", Position: 1}}}}}}}},
		}
	}
	if err := valid().Validate(); err != nil {
		t.Fatalf("valid course: %v", err)
	}
	for _, tc := range []struct {
		name   string
		change func(*Course)
	}{
		{"missing template", func(c *Course) { c.Template.Code = "" }},
		{"missing version", func(c *Course) { c.Version.Tag = "" }},
		{"module order", func(c *Course) { c.Modules[0].SequenceOrder = 2 }},
		{"lesson order", func(c *Course) { c.Modules[0].Lessons[0].SequenceOrder = 0 }},
		{"activity type", func(c *Course) { c.Modules[0].Lessons[0].Activities[0].Type = "" }},
		{"resource key", func(c *Course) { c.Modules[0].Lessons[0].Activities[0].Resources[0].URLStorageKey = "" }},
		{"resource size", func(c *Course) { c.Modules[0].Lessons[0].Activities[0].Resources[0].SizeBytes = -1 }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			course := valid()
			tc.change(&course)
			if err := course.Validate(); !errors.Is(err, ErrInvalid) {
				t.Fatalf("expected invalid, got %v", err)
			}
		})
	}
}

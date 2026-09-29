package domain

import "strings"

// Course is the catalog aggregate created as one unit with its initial version.
type Course struct {
	Template Template
	Version  Version
	Syllabus *Syllabus
	Modules  []Module
}

type Module struct {
	ID            string
	Title         string
	SequenceOrder int
	Lessons       []Lesson
}

type Lesson struct {
	ID            string
	Title         string
	SequenceOrder int
	Activities    []Activity
}

type Activity struct {
	ID            string
	Title         string
	Type          string
	IsRequired    bool
	SequenceOrder int
	Resources     []Resource
}

type Resource struct {
	ID            string
	Type          string
	Name          string
	URLStorageKey string
	MIMEType      string
	SizeBytes     int
	Position      int
	Metadata      map[string]any
}

func (c Course) Validate() error {
	if err := c.Template.Validate(); err != nil {
		return err
	}
	version := c.Version
	if version.TemplateID == "" {
		version.TemplateID = "pending-template"
	}
	if err := version.Validate(); err != nil {
		return err
	}
	if c.Syllabus != nil {
		syllabus := *c.Syllabus
		if syllabus.VersionID == "" {
			syllabus.VersionID = "pending-version"
		}
		if err := syllabus.Validate(); err != nil {
			return err
		}
	}
	for i, m := range c.Modules {
		if strings.TrimSpace(m.Title) == "" || m.SequenceOrder != i+1 {
			return ErrInvalid
		}
		for j, l := range m.Lessons {
			if strings.TrimSpace(l.Title) == "" || l.SequenceOrder != j+1 {
				return ErrInvalid
			}
			for k, a := range l.Activities {
				if strings.TrimSpace(a.Title) == "" || strings.TrimSpace(a.Type) == "" || a.SequenceOrder != k+1 {
					return ErrInvalid
				}
				for n, r := range a.Resources {
					if strings.TrimSpace(r.Type) == "" || strings.TrimSpace(r.Name) == "" || strings.TrimSpace(r.URLStorageKey) == "" || r.SizeBytes < 0 || r.Position != n+1 {
						return ErrInvalid
					}
				}
			}
		}
	}
	return nil
}

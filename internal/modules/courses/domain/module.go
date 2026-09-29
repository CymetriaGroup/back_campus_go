package domain

import (
	"errors"
	"strings"
)

var ErrInvalid = errors.New("invalid course data")
var ErrNotFound = errors.New("course item not found")
var ErrConflict = errors.New("course item already exists")

type Category struct{ ID, TenantID, Code, Name, Description, ParentID string }
type Template struct{ ID, Code, Title, Description string }
type Version struct {
	ID, TemplateID, Tag, Status string
	EstimatedHours              int
}
type Syllabus struct {
	ID, VersionID, Objectives, EntryProfile, ExitProfile, Methodology string
	DurationHours                                                     int
}

func (c Category) Validate() error {
	if strings.TrimSpace(c.TenantID) == "" || strings.TrimSpace(c.Code) == "" || strings.TrimSpace(c.Name) == "" || c.ID == c.ParentID {
		return ErrInvalid
	}
	return nil
}
func (t Template) Validate() error {
	if strings.TrimSpace(t.Code) == "" || strings.TrimSpace(t.Title) == "" {
		return ErrInvalid
	}
	return nil
}
func (v Version) Validate() error {
	if strings.TrimSpace(v.TemplateID) == "" || strings.TrimSpace(v.Tag) == "" || v.EstimatedHours < 0 || (v.Status != "DRAFT" && v.Status != "PUBLISHED" && v.Status != "ARCHIVED") {
		return ErrInvalid
	}
	return nil
}
func (s Syllabus) Validate() error {
	if strings.TrimSpace(s.VersionID) == "" || s.DurationHours < 0 {
		return ErrInvalid
	}
	return nil
}

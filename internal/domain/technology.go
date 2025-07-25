package domain

import (
	"time"
)

// Technology represents a technology skill or tool
type Technology struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	IconURL     string    `json:"icon_url,omitempty"`
	Level       Level     `json:"level"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Level represents proficiency level in a technology
type Level string

const (
	LevelBeginner     Level = "beginner"
	LevelIntermediate Level = "intermediate"
	LevelAdvanced     Level = "advanced"
	LevelExpert       Level = "expert"
)

// IsValid checks if the level is valid
func (l Level) IsValid() bool {
	switch l {
	case LevelBeginner, LevelIntermediate, LevelAdvanced, LevelExpert:
		return true
	default:
		return false
	}
}

// String returns the string representation of the level
func (l Level) String() string {
	return string(l)
}

// NewTechnology creates a new Technology instance
func NewTechnology(name, category string, level Level) *Technology {
	now := time.Now()
	return &Technology{
		Name:      name,
		Category:  category,
		Level:     level,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Validate checks if the Technology instance is valid
func (t *Technology) Validate() error {
	if t.Name == "" {
		return ErrInvalidInput("technology name cannot be empty")
	}
	if t.Category == "" {
		return ErrInvalidInput("technology category cannot be empty")
	}
	if !t.Level.IsValid() {
		return ErrInvalidInput("invalid technology level")
	}
	return nil
}

// UpdateLevel updates the proficiency level
func (t *Technology) UpdateLevel(level Level) error {
	if !level.IsValid() {
		return ErrInvalidInput("invalid technology level")
	}
	t.Level = level
	t.UpdatedAt = time.Now()
	return nil
}
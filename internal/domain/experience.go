package domain

import (
	"time"
)

// Experience represents a professional work experience
type Experience struct {
	ID           string        `json:"id"`
	CompanyName  string        `json:"company_name"`
	Position     string        `json:"position"`
	Description  string        `json:"description"`
	StartDate    time.Time     `json:"start_date"`
	EndDate      *time.Time    `json:"end_date,omitempty"` // nil for current position
	Location     string        `json:"location"`
	IsRemote     bool          `json:"is_remote"`
	Technologies []Technology  `json:"technologies"`
	Achievements []Achievement `json:"achievements"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// Achievement represents a specific achievement during an experience
type Achievement struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Impact      string    `json:"impact,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// NewExperience creates a new Experience instance
func NewExperience(companyName, position, description, location string, startDate time.Time, isRemote bool) *Experience {
	now := time.Now()
	return &Experience{
		CompanyName:  companyName,
		Position:     position,
		Description:  description,
		StartDate:    startDate,
		Location:     location,
		IsRemote:     isRemote,
		Technologies: make([]Technology, 0),
		Achievements: make([]Achievement, 0),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// NewAchievement creates a new Achievement instance
func NewAchievement(title, description, impact string) *Achievement {
	return &Achievement{
		Title:       title,
		Description: description,
		Impact:      impact,
		CreatedAt:   time.Now(),
	}
}

// Validate checks if the Experience instance is valid
func (e *Experience) Validate() error {
	if e.CompanyName == "" {
		return ErrInvalidInput("company name cannot be empty")
	}
	if e.Position == "" {
		return ErrInvalidInput("position cannot be empty")
	}
	if e.StartDate.IsZero() {
		return ErrInvalidInput("start date cannot be zero")
	}
	if e.EndDate != nil && e.EndDate.Before(e.StartDate) {
		return ErrInvalidInput("end date cannot be before start date")
	}
	return nil
}

// IsCurrent returns true if this is a current position
func (e *Experience) IsCurrent() bool {
	return e.EndDate == nil
}

// Duration calculates the duration of the experience
func (e *Experience) Duration() time.Duration {
	endDate := time.Now()
	if e.EndDate != nil {
		endDate = *e.EndDate
	}
	return endDate.Sub(e.StartDate)
}

// AddTechnology adds a technology to the experience
func (e *Experience) AddTechnology(tech Technology) {
	e.Technologies = append(e.Technologies, tech)
	e.UpdatedAt = time.Now()
}

// AddAchievement adds an achievement to the experience
func (e *Experience) AddAchievement(achievement Achievement) {
	e.Achievements = append(e.Achievements, achievement)
	e.UpdatedAt = time.Now()
}

// SetEndDate sets the end date for the experience
func (e *Experience) SetEndDate(endDate time.Time) error {
	if endDate.Before(e.StartDate) {
		return ErrInvalidInput("end date cannot be before start date")
	}
	e.EndDate = &endDate
	e.UpdatedAt = time.Now()
	return nil
}
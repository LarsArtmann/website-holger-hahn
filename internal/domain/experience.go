// Package domain provides core business domain entities and value objects for the portfolio website.
// It defines professional experience entities with quantified achievements and business metrics
// to showcase career progression and measurable impact across different roles.
package domain

import (
	"time"
)

// Experience represents a professional work experience.
type Experience struct {
	StartDate    time.Time     `json:"start_date"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	EndDate      *time.Time    `json:"end_date,omitempty"`
	ID           string        `json:"id"`
	CompanyName  string        `json:"company_name"`
	Position     string        `json:"position"`
	Description  string        `json:"description"`
	Location     string        `json:"location"`
	Technologies []Technology  `json:"technologies"`
	Achievements []Achievement `json:"achievements"`
	IsRemote     bool          `json:"is_remote"`
}

// Achievement represents a specific achievement during an experience.
type Achievement struct {
	CreatedAt   time.Time `json:"created_at"`
	Metrics     *Metrics  `json:"metrics,omitempty"`
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Impact      string    `json:"impact,omitempty"`
}

// Metrics represents quantified business impact metrics for achievements.
type Metrics struct {
	TestCoverage      *CoverageMetric     `json:"test_coverage,omitempty"`
	DeploymentTime    *TimeMetric         `json:"deployment_time,omitempty"`
	SystemReliability *ReliabilityMetric  `json:"system_reliability,omitempty"`
	Productivity      *ProductivityMetric `json:"productivity,omitempty"`
	CostSavings       *CostMetric         `json:"cost_savings,omitempty"`
}

// CoverageMetric represents test coverage improvements.
type CoverageMetric struct {
	Unit        string  `json:"unit"`
	Before      float64 `json:"before"`
	After       float64 `json:"after"`
	Improvement float64 `json:"improvement"`
}

// TimeMetric represents time-based improvements.
type TimeMetric struct {
	Unit        string  `json:"unit"`
	Before      int     `json:"before"`
	After       int     `json:"after"`
	Improvement float64 `json:"improvement"`
}

// ReliabilityMetric represents system reliability improvements.
type ReliabilityMetric struct {
	Unit      string  `json:"unit"`
	Uptime    float64 `json:"uptime"`
	MTBF      int     `json:"mtbf"`
	MTTR      int     `json:"mttr"`
	Incidents int     `json:"incidents"`
}

// ProductivityMetric represents team productivity improvements.
type ProductivityMetric struct {
	Unit                string  `json:"unit"`
	DeploymentFrequency int     `json:"deployment_frequency"`
	LeadTime            int     `json:"lead_time"`
	CycleTime           int     `json:"cycle_time"`
	Efficiency          float64 `json:"efficiency"`
}

// CostMetric represents cost savings and ROI.
type CostMetric struct {
	Unit           string  `json:"unit"`
	MonthlySavings float64 `json:"monthly_savings"`
	AnnualSavings  float64 `json:"annual_savings"`
	ROI            float64 `json:"roi"`
	PaybackPeriod  int     `json:"payback_period"`
}

// NewExperience creates a new Experience instance.
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

// NewAchievement creates a new Achievement instance.
func NewAchievement(title, description, impact string) *Achievement {
	return &Achievement{
		Title:       title,
		Description: description,
		Impact:      impact,
		CreatedAt:   time.Now(),
	}
}

// NewAchievementWithMetrics creates a new Achievement instance with metrics.
func NewAchievementWithMetrics(title, description, impact string, metrics *Metrics) *Achievement {
	return &Achievement{
		Title:       title,
		Description: description,
		Impact:      impact,
		Metrics:     metrics,
		CreatedAt:   time.Now(),
	}
}

// NewTestCoverageMetric creates a test coverage improvement metric.
func NewTestCoverageMetric(before, after float64) *CoverageMetric {
	improvement := after - before

	return &CoverageMetric{
		Before:      before,
		After:       after,
		Improvement: improvement,
		Unit:        "%",
	}
}

// NewDeploymentTimeMetric creates a deployment time improvement metric.
func NewDeploymentTimeMetric(beforeMinutes, afterMinutes int) *TimeMetric {
	improvement := float64(beforeMinutes-afterMinutes) / float64(beforeMinutes) * 100

	return &TimeMetric{
		Before:      beforeMinutes,
		After:       afterMinutes,
		Improvement: improvement,
		Unit:        "minutes",
	}
}

// NewReliabilityMetric creates a system reliability metric.
func NewReliabilityMetric(uptime float64, mtbf, mttr, incidents int) *ReliabilityMetric {
	return &ReliabilityMetric{
		Uptime:    uptime,
		MTBF:      mtbf,
		MTTR:      mttr,
		Incidents: incidents,
		Unit:      "%",
	}
}

// NewProductivityMetric creates a team productivity metric.
func NewProductivityMetric(deploymentFreq, leadTime, cycleTime int, efficiency float64) *ProductivityMetric {
	return &ProductivityMetric{
		DeploymentFrequency: deploymentFreq,
		LeadTime:            leadTime,
		CycleTime:           cycleTime,
		Efficiency:          efficiency,
		Unit:                "per week",
	}
}

// NewCostMetric creates a cost savings metric.
func NewCostMetric(monthlySavings, annualSavings, roi float64, paybackPeriod int) *CostMetric {
	return &CostMetric{
		MonthlySavings: monthlySavings,
		AnnualSavings:  annualSavings,
		ROI:            roi,
		PaybackPeriod:  paybackPeriod,
		Unit:           "USD",
	}
}

// Validate checks if the Experience instance is valid.
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

// IsCurrent returns true if this is a current position.
func (e *Experience) IsCurrent() bool {
	return e.EndDate == nil
}

// Duration calculates the duration of the experience.
func (e *Experience) Duration() time.Duration {
	endDate := time.Now()
	if e.EndDate != nil {
		endDate = *e.EndDate
	}

	return endDate.Sub(e.StartDate)
}

// AddTechnology adds a technology to the experience.
func (e *Experience) AddTechnology(tech Technology) {
	e.Technologies = append(e.Technologies, tech)
	e.UpdatedAt = time.Now()
}

// AddAchievement adds an achievement to the experience.
func (e *Experience) AddAchievement(achievement Achievement) {
	e.Achievements = append(e.Achievements, achievement)
	e.UpdatedAt = time.Now()
}

// GetID returns the experience ID.
func (e *Experience) GetID() string {
	return e.ID
}

// SetEndDate sets the end date for the experience.
func (e *Experience) SetEndDate(endDate time.Time) error {
	if endDate.Before(e.StartDate) {
		return ErrInvalidInput("end date cannot be before start date")
	}

	e.EndDate = &endDate
	e.UpdatedAt = time.Now()

	return nil
}

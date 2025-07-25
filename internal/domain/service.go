package domain

import (
	"time"
)

// Service represents a professional service offering
type Service struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Category     ServiceType   `json:"category"`
	Duration     string        `json:"duration,omitempty"`     // e.g., "3-6 months", "1 week"
	Pricing      *PricingInfo  `json:"pricing,omitempty"`
	Technologies []Technology  `json:"technologies"`
	Deliverables []Deliverable `json:"deliverables"`
	IsActive     bool          `json:"is_active"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// ServiceType represents the type/category of service
type ServiceType string

const (
	ServiceTypeConsulting   ServiceType = "consulting"
	ServiceTypeDevelopment  ServiceType = "development"
	ServiceTypeArchitecture ServiceType = "architecture"
	ServiceTypeAuditing     ServiceType = "auditing"
	ServiceTypeTraining     ServiceType = "training"
	ServiceTypeMentoring    ServiceType = "mentoring"
)

// PricingInfo represents pricing information for a service
type PricingInfo struct {
	Type        PricingType `json:"type"`
	Amount      float64     `json:"amount,omitempty"`
	Currency    string      `json:"currency,omitempty"`
	Description string      `json:"description,omitempty"`
}

// PricingType represents how the service is priced
type PricingType string

const (
	PricingTypeHourly    PricingType = "hourly"
	PricingTypeDaily     PricingType = "daily"
	PricingTypeProject   PricingType = "project"
	PricingTypeRetainer  PricingType = "retainer"
	PricingTypeCustom    PricingType = "custom"
)

// Deliverable represents what will be delivered as part of the service
type Deliverable struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Timeline    string    `json:"timeline,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// NewService creates a new Service instance
func NewService(name, description string, category ServiceType) *Service {
	now := time.Now()
	return &Service{
		Name:         name,
		Description:  description,
		Category:     category,
		Technologies: make([]Technology, 0),
		Deliverables: make([]Deliverable, 0),
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// NewDeliverable creates a new Deliverable instance
func NewDeliverable(name, description, timeline string) *Deliverable {
	return &Deliverable{
		Name:        name,
		Description: description,
		Timeline:    timeline,
		CreatedAt:   time.Now(),
	}
}

// IsValid checks if the service type is valid
func (st ServiceType) IsValid() bool {
	switch st {
	case ServiceTypeConsulting, ServiceTypeDevelopment, ServiceTypeArchitecture,
		 ServiceTypeAuditing, ServiceTypeTraining, ServiceTypeMentoring:
		return true
	default:
		return false
	}
}

// IsValid checks if the pricing type is valid
func (pt PricingType) IsValid() bool {
	switch pt {
	case PricingTypeHourly, PricingTypeDaily, PricingTypeProject,
		 PricingTypeRetainer, PricingTypeCustom:
		return true
	default:
		return false
	}
}

// Validate checks if the Service instance is valid
func (s *Service) Validate() error {
	if s.Name == "" {
		return ErrInvalidInput("service name cannot be empty")
	}
	if s.Description == "" {
		return ErrInvalidInput("service description cannot be empty")
	}
	if !s.Category.IsValid() {
		return ErrInvalidInput("invalid service category")
	}
	if s.Pricing != nil {
		if !s.Pricing.Type.IsValid() {
			return ErrInvalidInput("invalid pricing type")
		}
		if s.Pricing.Amount < 0 {
			return ErrInvalidInput("pricing amount cannot be negative")
		}
	}
	return nil
}

// AddTechnology adds a technology to the service
func (s *Service) AddTechnology(tech Technology) {
	s.Technologies = append(s.Technologies, tech)
	s.UpdatedAt = time.Now()
}

// AddDeliverable adds a deliverable to the service
func (s *Service) AddDeliverable(deliverable Deliverable) {
	s.Deliverables = append(s.Deliverables, deliverable)
	s.UpdatedAt = time.Now()
}

// SetPricing sets the pricing information for the service
func (s *Service) SetPricing(pricing PricingInfo) error {
	if !pricing.Type.IsValid() {
		return ErrInvalidInput("invalid pricing type")
	}
	if pricing.Amount < 0 {
		return ErrInvalidInput("pricing amount cannot be negative")
	}
	s.Pricing = &pricing
	s.UpdatedAt = time.Now()
	return nil
}

// Activate activates the service
func (s *Service) Activate() {
	s.IsActive = true
	s.UpdatedAt = time.Now()
}

// Deactivate deactivates the service
func (s *Service) Deactivate() {
	s.IsActive = false
	s.UpdatedAt = time.Now()
}
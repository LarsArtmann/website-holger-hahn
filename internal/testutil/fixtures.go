// Package testutil provides testing utilities, fixtures, and helpers for unit and integration tests.
// It includes common test data, mock implementations, and assertion helpers.
package testutil

import (
	"time"

	"holger-hahn-website/internal/domain"
)

// TestTimeNow provides a fixed time for consistent testing
var TestTimeNow = time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

// TechnologyFixtures provides common technology test data
type TechnologyFixtures struct{}

// NewTechnologyFixtures creates a new technology fixtures instance
func NewTechnologyFixtures() *TechnologyFixtures {
	return &TechnologyFixtures{}
}

// ValidTechnology returns a valid technology for testing
func (f *TechnologyFixtures) ValidTechnology() *domain.Technology {
	tech := domain.NewTechnology("Go", "Backend", domain.LevelExpert)
	tech.ID = "tech-001"
	tech.Description = "Server-side programming language"
	tech.IconURL = "https://example.com/go-icon.png"
	tech.CreatedAt = TestTimeNow
	tech.UpdatedAt = TestTimeNow
	return tech
}

// TechnologiesList returns a list of technologies for testing
func (f *TechnologyFixtures) TechnologiesList() []*domain.Technology {
	return []*domain.Technology{
		{
			ID:          "tech-001",
			Name:        "Go",
			Category:    "Backend",
			Level:       domain.LevelExpert,
			Description: "Server-side programming language",
			IconURL:     "https://example.com/go-icon.png",
			CreatedAt:   TestTimeNow,
			UpdatedAt:   TestTimeNow,
		},
		{
			ID:          "tech-002",
			Name:        "TypeScript",
			Category:    "Frontend",
			Level:       domain.LevelAdvanced,
			Description: "Typed JavaScript",
			IconURL:     "https://example.com/ts-icon.png",
			CreatedAt:   TestTimeNow,
			UpdatedAt:   TestTimeNow,
		},
		{
			ID:          "tech-003",
			Name:        "PostgreSQL",
			Category:    "Database",
			Level:       domain.LevelIntermediate,
			Description: "Relational database",
			IconURL:     "https://example.com/postgres-icon.png",
			CreatedAt:   TestTimeNow,
			UpdatedAt:   TestTimeNow,
		},
		{
			ID:          "tech-004",
			Name:        "Docker",
			Category:    "DevOps",
			Level:       domain.LevelAdvanced,
			Description: "Containerization platform",
			IconURL:     "https://example.com/docker-icon.png",
			CreatedAt:   TestTimeNow,
			UpdatedAt:   TestTimeNow,
		},
		{
			ID:          "tech-005",
			Name:        "Kubernetes",
			Category:    "DevOps",
			Level:       domain.LevelBeginner,
			Description: "Container orchestration",
			IconURL:     "https://example.com/k8s-icon.png",
			CreatedAt:   TestTimeNow,
			UpdatedAt:   TestTimeNow,
		},
	}
}

// ExpertTechnologies returns technologies with expert level
func (f *TechnologyFixtures) ExpertTechnologies() []*domain.Technology {
	all := f.TechnologiesList()
	var expert []*domain.Technology
	for _, tech := range all {
		if tech.Level == domain.LevelExpert {
			expert = append(expert, tech)
		}
	}
	return expert
}

// BackendTechnologies returns backend category technologies
func (f *TechnologyFixtures) BackendTechnologies() []*domain.Technology {
	all := f.TechnologiesList()
	var backend []*domain.Technology
	for _, tech := range all {
		if tech.Category == "Backend" {
			backend = append(backend, tech)
		}
	}
	return backend
}

// ExperienceFixtures provides common experience test data
type ExperienceFixtures struct{}

// NewExperienceFixtures creates a new experience fixtures instance
func NewExperienceFixtures() *ExperienceFixtures {
	return &ExperienceFixtures{}
}

// ValidExperience returns a valid experience for testing
func (f *ExperienceFixtures) ValidExperience() *domain.Experience {
	startDate := TestTimeNow.AddDate(-2, 0, 0) // 2 years ago
	exp := domain.NewExperience(
		"TechCorp Inc",
		"Senior Software Engineer",
		"Led development of cloud-native applications",
		"San Francisco, CA",
		startDate,
		false,
	)
	exp.ID = "exp-001"
	exp.CreatedAt = TestTimeNow
	exp.UpdatedAt = TestTimeNow
	return exp
}

// CurrentExperience returns a current experience (no end date)
func (f *ExperienceFixtures) CurrentExperience() *domain.Experience {
	startDate := TestTimeNow.AddDate(-1, 0, 0) // 1 year ago
	exp := domain.NewExperience(
		"InnovateLab",
		"Lead Developer",
		"Architecting scalable microservices",
		"Remote",
		startDate,
		true,
	)
	exp.ID = "exp-002"
	exp.CreatedAt = TestTimeNow
	exp.UpdatedAt = TestTimeNow
	return exp
}

// CompletedExperience returns a completed experience (with end date)
func (f *ExperienceFixtures) CompletedExperience() *domain.Experience {
	startDate := TestTimeNow.AddDate(-3, 0, 0) // 3 years ago
	endDate := TestTimeNow.AddDate(-2, 0, 0)   // 2 years ago
	exp := domain.NewExperience(
		"StartupXYZ",
		"Full Stack Developer",
		"Built MVP and scaled to production",
		"New York, NY",
		startDate,
		false,
	)
	exp.ID = "exp-003"
	exp.SetEndDate(endDate)
	exp.CreatedAt = TestTimeNow
	exp.UpdatedAt = TestTimeNow
	return exp
}

// ExperiencesList returns a list of experiences for testing
func (f *ExperienceFixtures) ExperiencesList() []*domain.Experience {
	return []*domain.Experience{
		f.CurrentExperience(),
		f.ValidExperience(),
		f.CompletedExperience(),
	}
}

// ExperienceWithTechnologies returns an experience with technologies added
func (f *ExperienceFixtures) ExperienceWithTechnologies() *domain.Experience {
	exp := f.ValidExperience()
	techFixtures := NewTechnologyFixtures()
	
	// Add some technologies
	exp.AddTechnology(*techFixtures.TechnologiesList()[0]) // Go
	exp.AddTechnology(*techFixtures.TechnologiesList()[1]) // TypeScript
	exp.AddTechnology(*techFixtures.TechnologiesList()[3]) // Docker
	
	return exp
}

// ExperienceWithAchievements returns an experience with achievements
func (f *ExperienceFixtures) ExperienceWithAchievements() *domain.Experience {
	exp := f.ValidExperience()
	
	// Add achievements with metrics
	testCoverage := domain.NewTestCoverageMetric(45.0, 85.0)
	deploymentTime := domain.NewDeploymentTimeMetric(120, 15)
	reliability := domain.NewReliabilityMetric(99.5, 720, 5, 2)
	
	metrics1 := &domain.Metrics{
		TestCoverage:      testCoverage,
		DeploymentTime:    deploymentTime,
		SystemReliability: reliability,
	}
	
	achievement1 := domain.NewAchievementWithMetrics(
		"Improved Test Coverage",
		"Implemented comprehensive testing strategy",
		"Reduced bugs in production by 60%",
		metrics1,
	)
	achievement1.ID = "ach-001"
	
	productivity := domain.NewProductivityMetric(5, 48, 24, 30.0)
	costSavings := domain.NewCostMetric(5000, 60000, 250.0, 6)
	
	metrics2 := &domain.Metrics{
		Productivity: productivity,
		CostSavings:  costSavings,
	}
	
	achievement2 := domain.NewAchievementWithMetrics(
		"Automated Deployment Pipeline",
		"Built CI/CD pipeline with automated testing",
		"Increased deployment frequency by 400%",
		metrics2,
	)
	achievement2.ID = "ach-002"
	
	exp.AddAchievement(*achievement1)
	exp.AddAchievement(*achievement2)
	
	return exp
}

// MetricsFixtures provides common metrics test data
type MetricsFixtures struct{}

// NewMetricsFixtures creates a new metrics fixtures instance
func NewMetricsFixtures() *MetricsFixtures {
	return &MetricsFixtures{}
}

// TestCoverageImprovement returns test coverage metrics
func (f *MetricsFixtures) TestCoverageImprovement() *domain.CoverageMetric {
	return domain.NewTestCoverageMetric(35.0, 80.0)
}

// DeploymentTimeImprovement returns deployment time metrics
func (f *MetricsFixtures) DeploymentTimeImprovement() *domain.TimeMetric {
	return domain.NewDeploymentTimeMetric(180, 10)
}

// SystemReliability returns reliability metrics
func (f *MetricsFixtures) SystemReliability() *domain.ReliabilityMetric {
	return domain.NewReliabilityMetric(99.9, 1440, 3, 1)
}

// TeamProductivity returns productivity metrics
func (f *MetricsFixtures) TeamProductivity() *domain.ProductivityMetric {
	return domain.NewProductivityMetric(8, 36, 18, 45.0)
}

// CostSavings returns cost savings metrics
func (f *MetricsFixtures) CostSavings() *domain.CostMetric {
	return domain.NewCostMetric(8000, 96000, 320.0, 4)
}

// CompleteMetrics returns a full metrics object
func (f *MetricsFixtures) CompleteMetrics() *domain.Metrics {
	return &domain.Metrics{
		TestCoverage:      f.TestCoverageImprovement(),
		DeploymentTime:    f.DeploymentTimeImprovement(),
		SystemReliability: f.SystemReliability(),
		Productivity:      f.TeamProductivity(),
		CostSavings:       f.CostSavings(),
	}
}
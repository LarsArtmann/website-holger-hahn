package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"holger-hahn-website/internal/domain"
	"holger-hahn-website/internal/testutil"
)

func TestExperienceService_CreateExperience(t *testing.T) {
	ctx := testutil.TestContext(t)
	mockExperienceRepo := testutil.NewMockExperienceRepository()
	mockTechRepo := testutil.NewMockTechnologyRepository()
	service := NewExperienceService(mockExperienceRepo, mockTechRepo)

	t.Run("successful creation", func(t *testing.T) {
		mockExperienceRepo.ClearCallLog()

		req := CreateExperienceRequest{
			CompanyName: "TechCorp",
			Position:    "Senior Developer",
			Description: "Led development of cloud-native applications",
			Location:    "San Francisco, CA",
			StartDate:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			IsRemote:    false,
		}

		exp, err := service.CreateExperience(ctx, req)

		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, exp)
		testutil.AssertEqual(t, "TechCorp", exp.CompanyName)
		testutil.AssertEqual(t, "Senior Developer", exp.Position)
		testutil.AssertEqual(t, "Led development of cloud-native applications", exp.Description)
		testutil.AssertEqual(t, "San Francisco, CA", exp.Location)
		testutil.AssertFalse(t, exp.IsRemote, "Should not be remote")
		testutil.AssertNotEqual(t, "", exp.ID)

		// Verify repository calls
		calls := mockExperienceRepo.GetCallLog()
		testutil.AssertLen(t, calls, 1)
		testutil.AssertTrue(t, calls[0][:6] == "Create", "Expected Create call")
	})

	t.Run("successful creation with end date", func(t *testing.T) {
		mockExperienceRepo.ClearCallLog()

		startDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

		req := CreateExperienceRequest{
			CompanyName: "OldCorp",
			Position:    "Developer",
			Description: "Maintained legacy systems",
			Location:    "New York, NY",
			StartDate:   startDate,
			EndDate:     &endDate,
			IsRemote:    true,
		}

		exp, err := service.CreateExperience(ctx, req)

		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, exp)
		testutil.AssertEqual(t, "OldCorp", exp.CompanyName)
		testutil.AssertTrue(t, exp.IsRemote, "Should be remote")
		testutil.AssertNotNil(t, exp.EndDate)
		testutil.AssertFalse(t, exp.IsCurrent(), "Should not be current")
	})

	t.Run("empty company name", func(t *testing.T) {
		req := CreateExperienceRequest{
			CompanyName: "",
			Position:    "Developer",
			Description: "Some description",
			StartDate:   time.Now(),
		}

		_, err := service.CreateExperience(ctx, req)

		testutil.AssertValidationError(t, err)
	})

	t.Run("whitespace-only company name", func(t *testing.T) {
		req := CreateExperienceRequest{
			CompanyName: "   ",
			Position:    "Developer",
			Description: "Some description",
			StartDate:   time.Now(),
		}

		_, err := service.CreateExperience(ctx, req)

		testutil.AssertValidationError(t, err)
	})

	t.Run("empty position", func(t *testing.T) {
		req := CreateExperienceRequest{
			CompanyName: "TechCorp",
			Position:    "",
			Description: "Some description",
			StartDate:   time.Now(),
		}

		_, err := service.CreateExperience(ctx, req)

		testutil.AssertValidationError(t, err)
	})

	t.Run("invalid end date before start date", func(t *testing.T) {
		startDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

		req := CreateExperienceRequest{
			CompanyName: "TechCorp",
			Position:    "Developer",
			Description: "Some description",
			StartDate:   startDate,
			EndDate:     &endDate,
		}

		_, err := service.CreateExperience(ctx, req)

		testutil.AssertValidationError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		mockExperienceRepo.SetErrorMode(true, "database connection failed")

		req := CreateExperienceRequest{
			CompanyName: "TechCorp",
			Position:    "Developer",
			Description: "Some description",
			StartDate:   time.Now(),
		}

		_, err := service.CreateExperience(ctx, req)

		testutil.AssertInternalError(t, err)
		mockExperienceRepo.SetErrorMode(false, "")
	})
}

func TestExperienceService_GetExperience(t *testing.T) {
	ctx := testutil.TestContext(t)
	mockExperienceRepo := testutil.NewMockExperienceRepository()
	mockTechRepo := testutil.NewMockTechnologyRepository()
	service := NewExperienceService(mockExperienceRepo, mockTechRepo)
	fixtures := testutil.NewExperienceFixtures()

	t.Run("successful retrieval", func(t *testing.T) {
		expected := fixtures.ValidExperience()
		mockExperienceRepo.PreloadExperiences([]*domain.Experience{expected})
		mockExperienceRepo.ClearCallLog()

		exp, err := service.GetExperience(ctx, expected.ID)

		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, exp)
		testutil.AssertEqual(t, expected.ID, exp.ID)
		testutil.AssertEqual(t, expected.CompanyName, exp.CompanyName)

		// Verify repository calls
		calls := mockExperienceRepo.GetCallLog()
		testutil.AssertLen(t, calls, 1)
		testutil.AssertContains(t, calls, "GetByID("+expected.ID+")")
	})

	t.Run("empty ID", func(t *testing.T) {
		_, err := service.GetExperience(ctx, "")

		testutil.AssertValidationError(t, err)
	})

	t.Run("experience not found", func(t *testing.T) {
		_, err := service.GetExperience(ctx, "non-existent-id")

		testutil.AssertNotFoundError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		mockExperienceRepo.SetErrorMode(true, "database connection failed")

		_, err := service.GetExperience(ctx, "any-id")

		testutil.AssertError(t, err)
		mockExperienceRepo.SetErrorMode(false, "")
	})
}

func TestExperienceService_ListExperiences(t *testing.T) {
	ctx := testutil.TestContext(t)
	mockExperienceRepo := testutil.NewMockExperienceRepository()
	mockTechRepo := testutil.NewMockTechnologyRepository()
	service := NewExperienceService(mockExperienceRepo, mockTechRepo)
	fixtures := testutil.NewExperienceFixtures()

	t.Run("list all experiences", func(t *testing.T) {
		expected := fixtures.ExperiencesList()
		mockExperienceRepo.PreloadExperiences(expected)
		mockExperienceRepo.ClearCallLog()

		experiences, err := service.ListExperiences(ctx, ExperienceFilter{})

		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, experiences)
		testutil.AssertLen(t, experiences, len(expected))

		// Verify repository calls
		calls := mockExperienceRepo.GetCallLog()
		testutil.AssertLen(t, calls, 1)
		testutil.AssertContains(t, calls, "List(filter)")
	})

	t.Run("list with company filter", func(t *testing.T) {
		expected := fixtures.ExperiencesList()
		mockExperienceRepo.PreloadExperiences(expected)
		mockExperienceRepo.ClearCallLog()

		companyName := "InnovateLab"
		filter := ExperienceFilter{CompanyName: &companyName}
		experiences, err := service.ListExperiences(ctx, filter)

		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, experiences)

		// All returned experiences should be from InnovateLab
		for _, exp := range experiences {
			testutil.AssertEqual(t, "InnovateLab", exp.CompanyName)
		}
	})

	t.Run("list with remote filter", func(t *testing.T) {
		expected := fixtures.ExperiencesList()
		mockExperienceRepo.PreloadExperiences(expected)
		mockExperienceRepo.ClearCallLog()

		isRemote := true
		filter := ExperienceFilter{IsRemote: &isRemote}
		experiences, err := service.ListExperiences(ctx, filter)

		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, experiences)

		// All returned experiences should be remote
		for _, exp := range experiences {
			testutil.AssertTrue(t, exp.IsRemote, "Should be remote")
		}
	})

	t.Run("list current experiences", func(t *testing.T) {
		expected := fixtures.ExperiencesList()
		mockExperienceRepo.PreloadExperiences(expected)
		mockExperienceRepo.ClearCallLog()

		isCurrent := true
		filter := ExperienceFilter{IsCurrent: &isCurrent}
		experiences, err := service.ListExperiences(ctx, filter)

		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, experiences)

		// All returned experiences should be current
		for _, exp := range experiences {
			testutil.AssertTrue(t, exp.IsCurrent(), "Should be current")
		}
	})

	t.Run("repository error", func(t *testing.T) {
		mockExperienceRepo.SetErrorMode(true, "database connection failed")

		_, err := service.ListExperiences(ctx, ExperienceFilter{})

		testutil.AssertInternalError(t, err)
		mockExperienceRepo.SetErrorMode(false, "")
	})
}

func TestExperienceService_GetCurrentExperiences(t *testing.T) {
	ctx := testutil.TestContext(t)
	mockExperienceRepo := testutil.NewMockExperienceRepository()
	mockTechRepo := testutil.NewMockTechnologyRepository()
	service := NewExperienceService(mockExperienceRepo, mockTechRepo)
	fixtures := testutil.NewExperienceFixtures()

	t.Run("successful retrieval", func(t *testing.T) {
		allExperiences := fixtures.ExperiencesList()
		mockExperienceRepo.PreloadExperiences(allExperiences)
		mockExperienceRepo.ClearCallLog()

		experiences, err := service.GetCurrentExperiences(ctx)

		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, experiences)

		// All returned experiences should be current
		for _, exp := range experiences {
			testutil.AssertTrue(t, exp.IsCurrent(), "Should be current")
		}

		// Verify repository calls
		calls := mockExperienceRepo.GetCallLog()
		testutil.AssertLen(t, calls, 1)
		testutil.AssertContains(t, calls, "GetCurrent()")
	})

	t.Run("repository error", func(t *testing.T) {
		mockExperienceRepo.SetErrorMode(true, "database connection failed")

		_, err := service.GetCurrentExperiences(ctx)

		testutil.AssertInternalError(t, err)
		mockExperienceRepo.SetErrorMode(false, "")
	})
}

func TestExperienceService_AddTechnologyToExperience(t *testing.T) {
	ctx := testutil.TestContext(t)
	mockExperienceRepo := testutil.NewMockExperienceRepository()
	mockTechRepo := testutil.NewMockTechnologyRepository()
	service := NewExperienceService(mockExperienceRepo, mockTechRepo)
	expFixtures := testutil.NewExperienceFixtures()
	techFixtures := testutil.NewTechnologyFixtures()

	t.Run("successful addition", func(t *testing.T) {
		experience := expFixtures.ValidExperience()
		technology := techFixtures.ValidTechnology()

		mockExperienceRepo.PreloadExperiences([]*domain.Experience{experience})
		mockTechRepo.PreloadTechnologies([]*domain.Technology{technology})
		mockExperienceRepo.ClearCallLog()
		mockTechRepo.ClearCallLog()

		err := service.AddTechnologyToExperience(ctx, experience.ID, technology.ID)

		testutil.AssertNoError(t, err)

		// Verify repository calls
		expCalls := mockExperienceRepo.GetCallLog()
		techCalls := mockTechRepo.GetCallLog()
		testutil.AssertTrue(t, len(expCalls) >= 2, "Should have GetByID and Update calls")
		testutil.AssertLen(t, techCalls, 1)
		testutil.AssertContains(t, techCalls, "GetByID("+technology.ID+")")
	})

	t.Run("experience not found", func(t *testing.T) {
		technology := techFixtures.ValidTechnology()
		mockTechRepo.PreloadTechnologies([]*domain.Technology{technology})

		err := service.AddTechnologyToExperience(ctx, "non-existent-exp", technology.ID)

		testutil.AssertNotFoundError(t, err)
	})

	t.Run("technology not found", func(t *testing.T) {
		experience := expFixtures.ValidExperience()
		mockExperienceRepo.PreloadExperiences([]*domain.Experience{experience})

		err := service.AddTechnologyToExperience(ctx, experience.ID, "non-existent-tech")

		testutil.AssertNotFoundError(t, err)
	})

	t.Run("technology already added", func(t *testing.T) {
		experience := expFixtures.ExperienceWithTechnologies()
		technology := experience.Technologies[0] // Use existing technology

		mockExperienceRepo.PreloadExperiences([]*domain.Experience{experience})
		mockTechRepo.PreloadTechnologies([]*domain.Technology{&technology})

		err := service.AddTechnologyToExperience(ctx, experience.ID, technology.ID)

		testutil.AssertConflictError(t, err)
	})
}

func TestExperienceService_AddAchievementToExperience(t *testing.T) {
	ctx := testutil.TestContext(t)
	mockExperienceRepo := testutil.NewMockExperienceRepository()
	mockTechRepo := testutil.NewMockTechnologyRepository()
	service := NewExperienceService(mockExperienceRepo, mockTechRepo)
	fixtures := testutil.NewExperienceFixtures()

	t.Run("successful addition without metrics", func(t *testing.T) {
		experience := fixtures.ValidExperience()
		mockExperienceRepo.PreloadExperiences([]*domain.Experience{experience})
		mockExperienceRepo.ClearCallLog()

		achievement := AchievementRequest{
			Title:       "Improved Performance",
			Description: "Optimized database queries",
			Impact:      "Reduced response time by 50%",
		}

		err := service.AddAchievementToExperience(ctx, experience.ID, achievement)

		testutil.AssertNoError(t, err)

		// Verify repository calls
		calls := mockExperienceRepo.GetCallLog()
		testutil.AssertTrue(t, len(calls) >= 2, "Should have GetByID and Update calls")
	})

	t.Run("successful addition with metrics", func(t *testing.T) {
		experience := fixtures.ValidExperience()
		mockExperienceRepo.PreloadExperiences([]*domain.Experience{experience})
		mockExperienceRepo.ClearCallLog()

		metricsFixtures := testutil.NewMetricsFixtures()
		metrics := metricsFixtures.CompleteMetrics()

		achievement := AchievementRequest{
			Title:       "Improved Test Coverage",
			Description: "Implemented comprehensive testing",
			Impact:      "Reduced production bugs",
			Metrics:     metrics,
		}

		err := service.AddAchievementToExperience(ctx, experience.ID, achievement)

		testutil.AssertNoError(t, err)
	})

	t.Run("experience not found", func(t *testing.T) {
		achievement := AchievementRequest{
			Title:       "Test Achievement",
			Description: "Test description",
		}

		err := service.AddAchievementToExperience(ctx, "non-existent-id", achievement)

		testutil.AssertNotFoundError(t, err)
	})
}

func TestExperienceService_EndExperience(t *testing.T) {
	ctx := testutil.TestContext(t)
	mockExperienceRepo := testutil.NewMockExperienceRepository()
	mockTechRepo := testutil.NewMockTechnologyRepository()
	service := NewExperienceService(mockExperienceRepo, mockTechRepo)
	fixtures := testutil.NewExperienceFixtures()

	t.Run("successful end", func(t *testing.T) {
		experience := fixtures.CurrentExperience() // Current experience (no end date)
		mockExperienceRepo.PreloadExperiences([]*domain.Experience{experience})
		mockExperienceRepo.ClearCallLog()

		endDate := time.Now()
		err := service.EndExperience(ctx, experience.ID, endDate)

		testutil.AssertNoError(t, err)

		// Verify repository calls
		calls := mockExperienceRepo.GetCallLog()
		testutil.AssertTrue(t, len(calls) >= 2, "Should have GetByID and Update calls")
	})

	t.Run("experience not found", func(t *testing.T) {
		endDate := time.Now()
		err := service.EndExperience(ctx, "non-existent-id", endDate)

		testutil.AssertNotFoundError(t, err)
	})

	t.Run("invalid end date before start date", func(t *testing.T) {
		experience := fixtures.ValidExperience()
		mockExperienceRepo.PreloadExperiences([]*domain.Experience{experience})

		// Try to set end date before start date
		endDate := experience.StartDate.AddDate(-1, 0, 0)
		err := service.EndExperience(ctx, experience.ID, endDate)

		testutil.AssertValidationError(t, err)
	})
}

func TestExperienceService_GetExperiencesByCompany(t *testing.T) {
	ctx := testutil.TestContext(t)
	mockExperienceRepo := testutil.NewMockExperienceRepository()
	mockTechRepo := testutil.NewMockTechnologyRepository()
	service := NewExperienceService(mockExperienceRepo, mockTechRepo)
	fixtures := testutil.NewExperienceFixtures()

	t.Run("successful retrieval", func(t *testing.T) {
		allExperiences := fixtures.ExperiencesList()
		mockExperienceRepo.PreloadExperiences(allExperiences)
		mockExperienceRepo.ClearCallLog()

		experiences, err := service.GetExperiencesByCompany(ctx, "InnovateLab")

		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, experiences)

		// All returned experiences should be from InnovateLab
		for _, exp := range experiences {
			testutil.AssertEqual(t, "InnovateLab", exp.CompanyName)
		}

		// Verify repository calls
		calls := mockExperienceRepo.GetCallLog()
		testutil.AssertLen(t, calls, 1)
		testutil.AssertContains(t, calls, "GetByCompany(InnovateLab)")
	})

	t.Run("empty company name", func(t *testing.T) {
		_, err := service.GetExperiencesByCompany(ctx, "")

		testutil.AssertValidationError(t, err)
	})
}

func TestExperienceService_GetTotalExperience(t *testing.T) {
	ctx := testutil.TestContext(t)
	mockExperienceRepo := testutil.NewMockExperienceRepository()
	mockTechRepo := testutil.NewMockTechnologyRepository()
	service := NewExperienceService(mockExperienceRepo, mockTechRepo)
	fixtures := testutil.NewExperienceFixtures()

	t.Run("successful calculation", func(t *testing.T) {
		allExperiences := fixtures.ExperiencesList()
		mockExperienceRepo.PreloadExperiences(allExperiences)
		mockExperienceRepo.ClearCallLog()

		duration, err := service.GetTotalExperience(ctx)

		testutil.AssertNoError(t, err)
		testutil.AssertTrue(t, duration > 0, "Total duration should be positive")

		// Verify repository calls
		calls := mockExperienceRepo.GetCallLog()
		testutil.AssertLen(t, calls, 1)
		testutil.AssertContains(t, calls, "List(filter)")
	})

	t.Run("repository error", func(t *testing.T) {
		mockExperienceRepo.SetErrorMode(true, "database connection failed")

		_, err := service.GetTotalExperience(ctx)

		testutil.AssertInternalError(t, err)
		mockExperienceRepo.SetErrorMode(false, "")
	})
}

// Benchmark tests
func BenchmarkExperienceService_CreateExperience(b *testing.B) {
	ctx := context.Background()
	mockExperienceRepo := testutil.NewMockExperienceRepository()
	mockTechRepo := testutil.NewMockTechnologyRepository()
	service := NewExperienceService(mockExperienceRepo, mockTechRepo)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := CreateExperienceRequest{
			CompanyName: fmt.Sprintf("Company-%d", i),
			Position:    "Developer",
			Description: "Test description",
			StartDate:   time.Now(),
		}
		_, _ = service.CreateExperience(ctx, req)
	}
}

func BenchmarkExperienceService_ListExperiences(b *testing.B) {
	ctx := context.Background()
	mockExperienceRepo := testutil.NewMockExperienceRepository()
	mockTechRepo := testutil.NewMockTechnologyRepository()
	service := NewExperienceService(mockExperienceRepo, mockTechRepo)

	// Preload some data
	fixtures := testutil.NewExperienceFixtures()
	mockExperienceRepo.PreloadExperiences(fixtures.ExperiencesList())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = service.ListExperiences(ctx, ExperienceFilter{})
	}
}

package service

import (
	"context"
	"fmt"
	"testing"

	"holger-hahn-website/internal/domain"
	"holger-hahn-website/internal/testutil"
)

func TestTechnologyService_CreateTechnology(t *testing.T) {
	ctx := testutil.TestContext(t)
	mockRepo := testutil.NewMockTechnologyRepository()
	service := NewTechnologyService(mockRepo)

	t.Run("successful creation", func(t *testing.T) {
		mockRepo.ClearCallLog()

		tech, err := service.CreateTechnology(ctx, "Go", "Backend", domain.LevelExpert)

		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, tech)
		testutil.AssertEqual(t, "Go", tech.Name)
		testutil.AssertEqual(t, "Backend", tech.Category)
		testutil.AssertEqual(t, domain.LevelExpert, tech.Level)
		testutil.AssertNotEqual(t, "", tech.ID)

		// Verify repository calls
		calls := mockRepo.GetCallLog()
		testutil.AssertLen(t, calls, 2)
		testutil.AssertContains(t, calls, "GetByName(Go)")
		testutil.AssertTrue(t, len(calls) >= 1 && calls[len(calls)-1][:6] == "Create", "Expected Create call")
	})

	t.Run("empty name", func(t *testing.T) {
		_, err := service.CreateTechnology(ctx, "", "Backend", domain.LevelExpert)

		testutil.AssertValidationError(t, err)
	})

	t.Run("whitespace-only name", func(t *testing.T) {
		_, err := service.CreateTechnology(ctx, "   ", "Backend", domain.LevelExpert)

		testutil.AssertValidationError(t, err)
	})

	t.Run("technology already exists", func(t *testing.T) {
		// Pre-load existing technology
		existing := testutil.NewTechnologyFixtures().ValidTechnology()
		mockRepo.PreloadTechnologies([]*domain.Technology{existing})

		_, err := service.CreateTechnology(ctx, existing.Name, "Backend", domain.LevelExpert)

		testutil.AssertConflictError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.SetErrorMode(true, "database connection failed")

		_, err := service.CreateTechnology(ctx, "NewTech", "Backend", domain.LevelExpert)

		testutil.AssertError(t, err)
		mockRepo.SetErrorMode(false, "")
	})

	t.Run("invalid level", func(t *testing.T) {
		_, err := service.CreateTechnology(ctx, "ValidName", "Backend", domain.Level("invalid"))

		testutil.AssertValidationError(t, err)
	})
}

func TestTechnologyService_GetTechnology(t *testing.T) {
	ctx := testutil.TestContext(t)
	mockRepo := testutil.NewMockTechnologyRepository()
	service := NewTechnologyService(mockRepo)
	fixtures := testutil.NewTechnologyFixtures()

	t.Run("successful retrieval", func(t *testing.T) {
		expected := fixtures.ValidTechnology()
		mockRepo.PreloadTechnologies([]*domain.Technology{expected})
		mockRepo.ClearCallLog()

		tech, err := service.GetTechnology(ctx, expected.ID)

		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, tech)
		testutil.AssertEqual(t, expected.ID, tech.ID)
		testutil.AssertEqual(t, expected.Name, tech.Name)

		// Verify repository calls
		calls := mockRepo.GetCallLog()
		testutil.AssertLen(t, calls, 1)
		testutil.AssertContains(t, calls, "GetByID("+expected.ID+")")
	})

	t.Run("empty ID", func(t *testing.T) {
		_, err := service.GetTechnology(ctx, "")

		testutil.AssertValidationError(t, err)
	})

	t.Run("technology not found", func(t *testing.T) {
		_, err := service.GetTechnology(ctx, "non-existent-id")

		testutil.AssertNotFoundError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.SetErrorMode(true, "database connection failed")

		_, err := service.GetTechnology(ctx, "any-id")

		testutil.AssertError(t, err)
		mockRepo.SetErrorMode(false, "")
	})
}

func TestTechnologyService_ListTechnologies(t *testing.T) {
	ctx := testutil.TestContext(t)
	mockRepo := testutil.NewMockTechnologyRepository()
	service := NewTechnologyService(mockRepo)
	fixtures := testutil.NewTechnologyFixtures()

	t.Run("list all technologies", func(t *testing.T) {
		expected := fixtures.TechnologiesList()
		mockRepo.PreloadTechnologies(expected)
		mockRepo.ClearCallLog()

		technologies, err := service.ListTechnologies(ctx, TechnologyFilter{})

		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, technologies)
		testutil.AssertLen(t, technologies, len(expected))

		// Verify repository calls
		calls := mockRepo.GetCallLog()
		testutil.AssertLen(t, calls, 1)
		testutil.AssertContains(t, calls, "List(filter)")
	})

	t.Run("list with category filter", func(t *testing.T) {
		expected := fixtures.TechnologiesList()
		mockRepo.PreloadTechnologies(expected)
		mockRepo.ClearCallLog()

		category := "Backend"
		filter := TechnologyFilter{Category: &category}
		technologies, err := service.ListTechnologies(ctx, filter)

		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, technologies)

		// All returned technologies should be Backend
		for _, tech := range technologies {
			testutil.AssertEqual(t, "Backend", tech.Category)
		}
	})

	t.Run("list with level filter", func(t *testing.T) {
		expected := fixtures.TechnologiesList()
		mockRepo.PreloadTechnologies(expected)
		mockRepo.ClearCallLog()

		level := domain.LevelExpert
		filter := TechnologyFilter{Level: &level}
		technologies, err := service.ListTechnologies(ctx, filter)

		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, technologies)

		// All returned technologies should be Expert level
		for _, tech := range technologies {
			testutil.AssertEqual(t, domain.LevelExpert, tech.Level)
		}
	})

	t.Run("list with pagination", func(t *testing.T) {
		expected := fixtures.TechnologiesList()
		mockRepo.PreloadTechnologies(expected)
		mockRepo.ClearCallLog()

		limit := 2
		offset := 1
		filter := TechnologyFilter{Limit: &limit, Offset: &offset}
		technologies, err := service.ListTechnologies(ctx, filter)

		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, technologies)
		testutil.AssertTrue(t, len(technologies) <= limit, "Should respect limit")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.SetErrorMode(true, "database connection failed")

		_, err := service.ListTechnologies(ctx, TechnologyFilter{})

		testutil.AssertInternalError(t, err)
		mockRepo.SetErrorMode(false, "")
	})
}

func TestTechnologyService_UpdateTechnology(t *testing.T) {
	ctx := testutil.TestContext(t)
	mockRepo := testutil.NewMockTechnologyRepository()
	service := NewTechnologyService(mockRepo)
	fixtures := testutil.NewTechnologyFixtures()

	t.Run("successful level update", func(t *testing.T) {
		existing := fixtures.ValidTechnology()
		mockRepo.PreloadTechnologies([]*domain.Technology{existing})
		mockRepo.ClearCallLog()

		newLevel := domain.LevelAdvanced
		updates := TechnologyUpdate{Level: &newLevel}

		tech, err := service.UpdateTechnology(ctx, existing.ID, updates)

		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, tech)
		testutil.AssertEqual(t, newLevel, tech.Level)

		// Verify repository calls
		calls := mockRepo.GetCallLog()
		testutil.AssertTrue(t, len(calls) >= 2, "Should have at least GetByID and Update calls")
	})

	t.Run("successful description update", func(t *testing.T) {
		existing := fixtures.ValidTechnology()
		mockRepo.PreloadTechnologies([]*domain.Technology{existing})
		mockRepo.ClearCallLog()

		newDescription := "Updated description"
		updates := TechnologyUpdate{Description: &newDescription}

		tech, err := service.UpdateTechnology(ctx, existing.ID, updates)

		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, tech)
		testutil.AssertEqual(t, newDescription, tech.Description)
	})

	t.Run("technology not found", func(t *testing.T) {
		updates := TechnologyUpdate{}

		_, err := service.UpdateTechnology(ctx, "non-existent-id", updates)

		testutil.AssertNotFoundError(t, err)
	})

	t.Run("invalid level update", func(t *testing.T) {
		existing := fixtures.ValidTechnology()
		mockRepo.PreloadTechnologies([]*domain.Technology{existing})

		invalidLevel := domain.Level("invalid")
		updates := TechnologyUpdate{Level: &invalidLevel}

		_, err := service.UpdateTechnology(ctx, existing.ID, updates)

		testutil.AssertValidationError(t, err)
	})
}

func TestTechnologyService_DeleteTechnology(t *testing.T) {
	ctx := testutil.TestContext(t)
	mockRepo := testutil.NewMockTechnologyRepository()
	service := NewTechnologyService(mockRepo)
	fixtures := testutil.NewTechnologyFixtures()

	t.Run("successful deletion", func(t *testing.T) {
		existing := fixtures.ValidTechnology()
		mockRepo.PreloadTechnologies([]*domain.Technology{existing})
		mockRepo.ClearCallLog()

		err := service.DeleteTechnology(ctx, existing.ID)

		testutil.AssertNoError(t, err)

		// Verify repository calls
		calls := mockRepo.GetCallLog()
		testutil.AssertTrue(t, len(calls) >= 2, "Should have GetByID and Delete calls")
		testutil.AssertTrue(t, calls[len(calls)-1] == "Delete("+existing.ID+")", "Should end with Delete call")
	})

	t.Run("empty ID", func(t *testing.T) {
		err := service.DeleteTechnology(ctx, "")

		testutil.AssertValidationError(t, err)
	})

	t.Run("technology not found", func(t *testing.T) {
		err := service.DeleteTechnology(ctx, "non-existent-id")

		testutil.AssertNotFoundError(t, err)
	})
}

func TestTechnologyService_GetTechnologiesByCategory(t *testing.T) {
	ctx := testutil.TestContext(t)
	mockRepo := testutil.NewMockTechnologyRepository()
	service := NewTechnologyService(mockRepo)
	fixtures := testutil.NewTechnologyFixtures()

	t.Run("successful retrieval", func(t *testing.T) {
		mockRepo.PreloadTechnologies(fixtures.TechnologiesList())
		mockRepo.ClearCallLog()

		technologies, err := service.GetTechnologiesByCategory(ctx, "Backend")

		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, technologies)

		// All returned technologies should be Backend
		for _, tech := range technologies {
			testutil.AssertEqual(t, "Backend", tech.Category)
		}

		// Verify repository calls
		calls := mockRepo.GetCallLog()
		testutil.AssertLen(t, calls, 1)
		testutil.AssertContains(t, calls, "GetByCategory(Backend)")
	})

	t.Run("empty category", func(t *testing.T) {
		_, err := service.GetTechnologiesByCategory(ctx, "")

		testutil.AssertValidationError(t, err)
	})
}

func TestTechnologyService_GetTechnologiesByLevel(t *testing.T) {
	ctx := testutil.TestContext(t)
	mockRepo := testutil.NewMockTechnologyRepository()
	service := NewTechnologyService(mockRepo)
	fixtures := testutil.NewTechnologyFixtures()

	t.Run("successful retrieval", func(t *testing.T) {
		mockRepo.PreloadTechnologies(fixtures.TechnologiesList())
		mockRepo.ClearCallLog()

		technologies, err := service.GetTechnologiesByLevel(ctx, domain.LevelExpert)

		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, technologies)

		// All returned technologies should be Expert level
		for _, tech := range technologies {
			testutil.AssertEqual(t, domain.LevelExpert, tech.Level)
		}

		// Verify repository calls
		calls := mockRepo.GetCallLog()
		testutil.AssertLen(t, calls, 1)
		testutil.AssertContains(t, calls, "GetByLevel(expert)")
	})

	t.Run("invalid level", func(t *testing.T) {
		_, err := service.GetTechnologiesByLevel(ctx, domain.Level("invalid"))

		testutil.AssertValidationError(t, err)
	})
}

// Benchmark tests
func BenchmarkTechnologyService_CreateTechnology(b *testing.B) {
	ctx := context.Background()
	mockRepo := testutil.NewMockTechnologyRepository()
	service := NewTechnologyService(mockRepo)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		name := fmt.Sprintf("Tech-%d", i)
		_, _ = service.CreateTechnology(ctx, name, "Backend", domain.LevelExpert)
	}
}

func BenchmarkTechnologyService_ListTechnologies(b *testing.B) {
	ctx := context.Background()
	mockRepo := testutil.NewMockTechnologyRepository()
	service := NewTechnologyService(mockRepo)

	// Preload some data
	fixtures := testutil.NewTechnologyFixtures()
	mockRepo.PreloadTechnologies(fixtures.TechnologiesList())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = service.ListTechnologies(ctx, TechnologyFilter{})
	}
}

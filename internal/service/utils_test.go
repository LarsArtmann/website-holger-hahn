package service

import (
	"errors"
	"strings"
	"testing"

	"holger-hahn-website/internal/testutil"
)

func TestGenerateID(t *testing.T) {
	t.Run("generates non-empty ID", func(t *testing.T) {
		id := generateID()

		testutil.AssertNotEqual(t, "", id)
		testutil.AssertTrue(t, len(id) > 10, "ID should be reasonably long")
	})

	t.Run("generates unique IDs", func(t *testing.T) {
		ids := make(map[string]bool)

		// Generate multiple IDs and check for uniqueness
		for range 100 {
			id := generateID()
			testutil.AssertFalse(t, ids[id], "ID should be unique: "+id)
			ids[id] = true
		}
	})

	t.Run("ID format contains timestamp and random part", func(t *testing.T) {
		id := generateID()

		// ID should contain a hyphen separating timestamp and random part
		parts := strings.Split(id, "-")
		testutil.AssertTrue(t, len(parts) >= 2, "ID should have timestamp-random format")

		// First part should be numeric (timestamp)
		timestamp := parts[0]
		testutil.AssertTrue(t, len(timestamp) > 0, "Timestamp part should not be empty")

		// Should have a random part
		randomPart := parts[1]
		testutil.AssertTrue(t, len(randomPart) > 0, "Random part should not be empty")
	})

	t.Run("handles error case gracefully", func(t *testing.T) {
		// Test that even if random generation fails, we get a valid ID
		// (though our implementation should handle this internally)
		id := generateID()

		// Should not contain "error" unless there's actually an error
		testutil.AssertFalse(t, strings.Contains(id, "error"), "Should not contain error in normal case")
	})
}

func TestServiceError(t *testing.T) {
	t.Run("creates service error with context", func(t *testing.T) {
		originalErr := errors.New("database connection failed")
		operation := "CreateTechnology"

		serviceErr := &ServiceError{
			Err:       originalErr,
			Operation: operation,
		}

		testutil.AssertEqual(t, originalErr, serviceErr.Err)
		testutil.AssertEqual(t, operation, serviceErr.Operation)
	})

	t.Run("Error method returns formatted message", func(t *testing.T) {
		originalErr := errors.New("database connection failed")
		operation := "UpdateExperience"

		serviceErr := &ServiceError{
			Err:       originalErr,
			Operation: operation,
		}

		expected := "service error in UpdateExperience: database connection failed"
		testutil.AssertEqual(t, expected, serviceErr.Error())
	})

	t.Run("Unwrap returns original error", func(t *testing.T) {
		originalErr := errors.New("original error")
		serviceErr := &ServiceError{
			Err:       originalErr,
			Operation: "TestOperation",
		}

		unwrapped := serviceErr.Unwrap()
		testutil.AssertEqual(t, originalErr, unwrapped)
	})
}

func TestWrapError(t *testing.T) {
	t.Run("wraps error with operation context", func(t *testing.T) {
		originalErr := errors.New("repository error")
		operation := "DeleteService"

		wrappedErr := WrapError(operation, originalErr)

		testutil.AssertNotNil(t, wrappedErr)

		serviceErr, ok := wrappedErr.(*ServiceError)
		testutil.AssertTrue(t, ok, "Should be a ServiceError")
		testutil.AssertEqual(t, originalErr, serviceErr.Err)
		testutil.AssertEqual(t, operation, serviceErr.Operation)
	})

	t.Run("returns nil for nil error", func(t *testing.T) {
		wrappedErr := WrapError("SomeOperation", nil)

		testutil.AssertNil(t, wrappedErr)
	})

	t.Run("wrapped error can be unwrapped", func(t *testing.T) {
		originalErr := errors.New("original error")
		wrappedErr := WrapError("TestOp", originalErr)

		unwrapped := errors.Unwrap(wrappedErr)
		testutil.AssertEqual(t, originalErr, unwrapped)
	})
}

func TestPagination(t *testing.T) {
	t.Run("NewPagination with valid values", func(t *testing.T) {
		pagination := NewPagination(20, 40)

		testutil.AssertEqual(t, 20, pagination.Limit)
		testutil.AssertEqual(t, 40, pagination.Offset)
		testutil.AssertEqual(t, 0, pagination.Total) // Total is not set by constructor
	})

	t.Run("NewPagination with zero limit uses default", func(t *testing.T) {
		pagination := NewPagination(0, 10)

		testutil.AssertEqual(t, 10, pagination.Limit) // Default limit
		testutil.AssertEqual(t, 10, pagination.Offset)
	})

	t.Run("NewPagination with negative limit uses default", func(t *testing.T) {
		pagination := NewPagination(-5, 5)

		testutil.AssertEqual(t, 10, pagination.Limit) // Default limit
		testutil.AssertEqual(t, 5, pagination.Offset)
	})

	t.Run("NewPagination caps limit at maximum", func(t *testing.T) {
		pagination := NewPagination(150, 0)

		testutil.AssertEqual(t, 100, pagination.Limit) // Max limit
		testutil.AssertEqual(t, 0, pagination.Offset)
	})

	t.Run("NewPagination with negative offset uses zero", func(t *testing.T) {
		pagination := NewPagination(20, -10)

		testutil.AssertEqual(t, 20, pagination.Limit)
		testutil.AssertEqual(t, 0, pagination.Offset) // Corrected to zero
	})

	t.Run("pagination values within expected ranges", func(t *testing.T) {
		testCases := []struct {
			inputLimit     int
			inputOffset    int
			expectedLimit  int
			expectedOffset int
		}{
			{1, 0, 1, 0},       // Minimum valid limit
			{100, 50, 100, 50}, // Maximum valid limit
			{50, 100, 50, 100}, // Normal case
			{0, 0, 10, 0},      // Default limit
			{-1, -1, 10, 0},    // Both negative
			{200, 0, 100, 0},   // Over max limit
		}

		for _, tc := range testCases {
			pagination := NewPagination(tc.inputLimit, tc.inputOffset)
			testutil.AssertEqual(t, tc.expectedLimit, pagination.Limit)
			testutil.AssertEqual(t, tc.expectedOffset, pagination.Offset)
		}
	})
}

func TestOrdering(t *testing.T) {
	t.Run("NewOrdering with valid values", func(t *testing.T) {
		ordering := NewOrdering("name", "asc")

		testutil.AssertEqual(t, "name", ordering.OrderBy)
		testutil.AssertEqual(t, "asc", ordering.OrderDir)
	})

	t.Run("NewOrdering with descending direction", func(t *testing.T) {
		ordering := NewOrdering("created_at", "desc")

		testutil.AssertEqual(t, "created_at", ordering.OrderBy)
		testutil.AssertEqual(t, "desc", ordering.OrderDir)
	})

	t.Run("NewOrdering with invalid direction defaults to asc", func(t *testing.T) {
		ordering := NewOrdering("name", "invalid")

		testutil.AssertEqual(t, "name", ordering.OrderBy)
		testutil.AssertEqual(t, "asc", ordering.OrderDir) // Default
	})

	t.Run("NewOrdering with empty direction defaults to asc", func(t *testing.T) {
		ordering := NewOrdering("name", "")

		testutil.AssertEqual(t, "name", ordering.OrderBy)
		testutil.AssertEqual(t, "asc", ordering.OrderDir) // Default
	})

	t.Run("NewOrdering preserves orderBy field", func(t *testing.T) {
		testFields := []string{"id", "name", "created_at", "updated_at", "category"}

		for _, field := range testFields {
			ordering := NewOrdering(field, "desc")
			testutil.AssertEqual(t, field, ordering.OrderBy)
			testutil.AssertEqual(t, "desc", ordering.OrderDir)
		}
	})
}

func TestValidateOrderBy(t *testing.T) {
	t.Run("valid technology fields", func(t *testing.T) {
		validFields := []string{"name", "category", "level", "created_at", "updated_at"}

		for _, field := range validFields {
			result := ValidateOrderBy("technology", field)
			testutil.AssertTrue(t, result, "Field should be valid: "+field)
		}
	})

	t.Run("valid experience fields", func(t *testing.T) {
		validFields := []string{"start_date", "end_date", "company_name", "position", "created_at", "updated_at"}

		for _, field := range validFields {
			result := ValidateOrderBy("experience", field)
			testutil.AssertTrue(t, result, "Field should be valid: "+field)
		}
	})

	t.Run("valid service fields", func(t *testing.T) {
		validFields := []string{"name", "category", "created_at", "updated_at"}

		for _, field := range validFields {
			result := ValidateOrderBy("service", field)
			testutil.AssertTrue(t, result, "Field should be valid: "+field)
		}
	})

	t.Run("invalid fields for technology", func(t *testing.T) {
		invalidFields := []string{"invalid_field", "description", "icon_url", "price"}

		for _, field := range invalidFields {
			result := ValidateOrderBy("technology", field)
			testutil.AssertFalse(t, result, "Field should be invalid: "+field)
		}
	})

	t.Run("invalid fields for experience", func(t *testing.T) {
		invalidFields := []string{"invalid_field", "description", "location", "salary"}

		for _, field := range invalidFields {
			result := ValidateOrderBy("experience", field)
			testutil.AssertFalse(t, result, "Field should be invalid: "+field)
		}
	})

	t.Run("invalid fields for service", func(t *testing.T) {
		invalidFields := []string{"invalid_field", "description", "price", "duration"}

		for _, field := range invalidFields {
			result := ValidateOrderBy("service", field)
			testutil.AssertFalse(t, result, "Field should be invalid: "+field)
		}
	})

	t.Run("invalid entity type", func(t *testing.T) {
		result := ValidateOrderBy("invalid_entity", "name")
		testutil.AssertFalse(t, result, "Invalid entity type should return false")
	})

	t.Run("empty entity type", func(t *testing.T) {
		result := ValidateOrderBy("", "name")
		testutil.AssertFalse(t, result, "Empty entity type should return false")
	})

	t.Run("empty field name", func(t *testing.T) {
		result := ValidateOrderBy("technology", "")
		testutil.AssertFalse(t, result, "Empty field name should return false")
	})

	t.Run("case sensitivity", func(t *testing.T) {
		// Fields should be case-sensitive
		result := ValidateOrderBy("technology", "Name") // Capital N
		testutil.AssertFalse(t, result, "Field names should be case-sensitive")

		result = ValidateOrderBy("Technology", "name") // Capital T
		testutil.AssertFalse(t, result, "Entity names should be case-sensitive")
	})
}

// Benchmark tests.
func BenchmarkGenerateID(b *testing.B) {
	for range b.N {
		_ = generateID()
	}
}

func BenchmarkWrapError(b *testing.B) {
	err := errors.New("test error")

	b.ResetTimer()

	for range b.N {
		_ = WrapError("TestOperation", err)
	}
}

func BenchmarkNewPagination(b *testing.B) {
	for i := range b.N {
		_ = NewPagination(20, i%1000)
	}
}

func BenchmarkValidateOrderBy(b *testing.B) {
	for range b.N {
		_ = ValidateOrderBy("technology", "name")
	}
}

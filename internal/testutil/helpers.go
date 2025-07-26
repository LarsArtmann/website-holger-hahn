package testutil

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"holger-hahn-website/internal/domain"
)

// AssertNoError fails the test if err is not nil
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

// AssertError fails the test if err is nil
func AssertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

// AssertEqual fails the test if expected != actual
func AssertEqual[T comparable](t *testing.T, expected, actual T) {
	t.Helper()
	if expected != actual {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}

// AssertNotEqual fails the test if expected == actual
func AssertNotEqual[T comparable](t *testing.T, expected, actual T) {
	t.Helper()
	if expected == actual {
		t.Fatalf("Expected values to be different, but both were %v", actual)
	}
}

// AssertNil fails the test if value is not nil
func AssertNil(t *testing.T, value interface{}) {
	t.Helper()
	if value != nil {
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Ptr && !rv.IsNil() {
			t.Fatalf("Expected nil, got %v", value)
		} else if rv.Kind() != reflect.Ptr && rv.IsValid() {
			t.Fatalf("Expected nil, got %v", value)
		}
	}
}

// AssertNotNil fails the test if value is nil
func AssertNotNil(t *testing.T, value interface{}) {
	t.Helper()
	if value == nil {
		t.Fatal("Expected non-nil value, got nil")
	}
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		t.Fatal("Expected non-nil pointer, got nil")
	}
}

// AssertTrue fails the test if condition is false
func AssertTrue(t *testing.T, condition bool, message string) {
	t.Helper()
	if !condition {
		t.Fatalf("Expected true: %s", message)
	}
}

// AssertFalse fails the test if condition is true
func AssertFalse(t *testing.T, condition bool, message string) {
	t.Helper()
	if condition {
		t.Fatalf("Expected false: %s", message)
	}
}

// AssertContains fails the test if slice doesn't contain item
func AssertContains[T comparable](t *testing.T, slice []T, item T) {
	t.Helper()
	for _, v := range slice {
		if v == item {
			return
		}
	}
	t.Fatalf("Expected slice to contain %v, but it didn't", item)
}

// AssertNotContains fails the test if slice contains item
func AssertNotContains[T comparable](t *testing.T, slice []T, item T) {
	t.Helper()
	for _, v := range slice {
		if v == item {
			t.Fatalf("Expected slice not to contain %v, but it did", item)
		}
	}
}

// AssertLen fails the test if length doesn't match expected
func AssertLen[T any](t *testing.T, slice []T, expectedLen int) {
	t.Helper()
	if len(slice) != expectedLen {
		t.Fatalf("Expected length %d, got %d", expectedLen, len(slice))
	}
}

// AssertTimeBetween fails the test if time is not between start and end
func AssertTimeBetween(t *testing.T, timeToCheck, start, end time.Time) {
	t.Helper()
	if timeToCheck.Before(start) || timeToCheck.After(end) {
		t.Fatalf("Expected time %v to be between %v and %v", timeToCheck, start, end)
	}
}

// AssertValidationError fails if error is not a validation error
func AssertValidationError(t *testing.T, err error) {
	t.Helper()
	AssertError(t, err)
	if !domain.IsValidationError(err) {
		t.Fatalf("Expected validation error, got: %v", err)
	}
}

// AssertNotFoundError fails if error is not a not found error
func AssertNotFoundError(t *testing.T, err error) {
	t.Helper()
	AssertError(t, err)
	if !domain.IsNotFoundError(err) {
		t.Fatalf("Expected not found error, got: %v", err)
	}
}

// AssertConflictError fails if error is not a conflict error
func AssertConflictError(t *testing.T, err error) {
	t.Helper()
	AssertError(t, err)
	if !domain.IsConflictError(err) {
		t.Fatalf("Expected conflict error, got: %v", err)
	}
}

// AssertInternalError fails if error is not an internal error
func AssertInternalError(t *testing.T, err error) {
	t.Helper()
	AssertError(t, err)
	if !domain.IsInternalError(err) {
		t.Fatalf("Expected internal error, got: %v", err)
	}
}

// WithEnv temporarily sets environment variables for testing
func WithEnv(t *testing.T, env map[string]string, test func()) {
	t.Helper()

	// Store original values
	original := make(map[string]string)
	for key := range env {
		original[key] = os.Getenv(key)
	}

	// Set test values
	for key, value := range env {
		os.Setenv(key, value)
	}

	// Restore original values after test
	defer func() {
		for key, value := range original {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	test()
}

// TestContext returns a context with timeout for testing
func TestContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	return ctx
}

// MockTime sets a fixed time for testing
func MockTime(t *testing.T, mockTime time.Time) func() {
	t.Helper()

	// This is a simple approach - in a real app you might use dependency injection
	// for time.Now() or a time interface
	return func() {
		// Restore function would go here if we had time mocking
	}
}

// AssertEventuallyTrue retries a condition with backoff until it's true or timeout
func AssertEventuallyTrue(t *testing.T, condition func() bool, timeout time.Duration, message string) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	interval := 10 * time.Millisecond

	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(interval)
		interval = time.Duration(float64(interval) * 1.5) // exponential backoff
		if interval > 100*time.Millisecond {
			interval = 100 * time.Millisecond
		}
	}

	t.Fatalf("Condition never became true within %v: %s", timeout, message)
}

// FakeTime represents a controllable time for testing
type FakeTime struct {
	current time.Time
}

// NewFakeTime creates a new fake time starting at the given time
func NewFakeTime(startTime time.Time) *FakeTime {
	return &FakeTime{current: startTime}
}

// Now returns the current fake time
func (f *FakeTime) Now() time.Time {
	return f.current
}

// Advance advances the fake time by the given duration
func (f *FakeTime) Advance(d time.Duration) {
	f.current = f.current.Add(d)
}

// Set sets the fake time to the given time
func (f *FakeTime) Set(t time.Time) {
	f.current = t
}

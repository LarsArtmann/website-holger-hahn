// Package utils provides utility functions and demonstrates usage of new libraries.
// This file serves as a placeholder to ensure important libraries stay in go.mod
// and provides examples for future integration.
package utils

import (
	"context"
	"errors"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Example functions demonstrating library usage to keep them in go.mod

// FilterStrings demonstrates samber/lo usage for functional programming
func FilterStrings(items []string, predicate func(string) bool) []string {
	return lo.Filter(items, func(item string, _ int) bool {
		return predicate(item)
	})
}

// MapStringsToInts demonstrates samber/lo mapping functionality
func MapStringsToInts(items []string, mapper func(string) int) []int {
	return lo.Map(items, func(item string, _ int) int {
		return mapper(item)
	})
}

// SafeOperation demonstrates samber/mo Result type for error handling
func SafeOperation(input string) mo.Result[string] {
	if input == "" {
		return mo.Err[string](errors.New("input cannot be empty"))
	}
	return mo.Ok("processed: " + input)
}

// GetConfigValue demonstrates spf13/viper usage for configuration
func GetConfigValue(key string, defaultValue interface{}) interface{} {
	viper.SetDefault(key, defaultValue)
	return viper.Get(key)
}

// TraceOperation demonstrates OpenTelemetry tracing
func TraceOperation(ctx context.Context, operationName string) (context.Context, trace.Span) {
	tracer := otel.Tracer("holger-hahn-website")
	return tracer.Start(ctx, operationName)
}

// Note: These functions ensure the libraries remain in go.mod and provide
// examples for future integration throughout the codebase.

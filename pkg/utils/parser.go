package utils

import (
	"strconv"
	"time"
)

// DerefOrDefault dereferences a pointer of any type.
// If the pointer is nil, it returns the provided default value.
func DerefOrDefault[T any](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// DerefString safely dereferences a *string, returning "" if it's nil.
func DerefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// parseFloat64 safely dereferences and parses a *string to a float64.
// Returns 0.0 if the pointer is nil or parsing fails.
func ParseFloat64(s *string) float64 {
	if s == nil {
		return 0.0
	}
	val, err := strconv.ParseFloat(*s, 64)
	if err != nil {
		return 0.0 // Atau nilai default lain yang Anda inginkan
	}
	return val
}

// ParseInt safely dereferences and parses a *string to an int.
// Returns 0 if the pointer is nil or parsing fails.
func ParseInt(s *string) int {
	if s == nil {
		return 0
	}
	val, err := strconv.Atoi(*s)
	if err != nil {
		return 0 // Atau nilai default lain, misal 1 untuk 'page'
	}
	return val
}

// ParseStringPointer simply returns the pointer itself if not nil.
// It's useful for consistency in function calls.
func ParseStringPointer(s *string) *string {
	if s == nil {
		return nil
	}
	return s
}

// ParseIntPointer parses a string pointer to an integer pointer.
// Returns nil if the input is nil or if parsing fails.
func ParseIntPointer(s *string) *int {
	if s == nil {
		return nil
	}
	val, err := strconv.Atoi(*s)
	if err != nil {
		return nil
	}
	return &val
}

// ParseFloat64Pointer parses a string pointer to a float64 pointer.
// Returns nil if the input is nil or if parsing fails.
func ParseFloat64Pointer(s *string) *float64 {
	if s == nil {
		return nil
	}
	val, err := strconv.ParseFloat(*s, 64)
	if err != nil {
		return nil
	}
	return &val
}

// ParseBoolPointer parses a string pointer to a boolean pointer.
// Returns nil if the input is nil or if parsing fails.
func ParseBoolPointer(s *string) *bool {
	if s == nil {
		return nil
	}
	val, err := strconv.ParseBool(*s)
	if err != nil {
		return nil
	}
	return &val
}

// ParseTimePointer parses a string pointer to a time.Time pointer using a specific layout.
// Returns nil if the input is nil or if parsing fails.
func ParseTimePointer(s *string, layout string) *time.Time {
	if s == nil {
		return nil
	}
	val, err := time.Parse(layout, *s)
	if err != nil {
		return nil
	}
	return &val
}

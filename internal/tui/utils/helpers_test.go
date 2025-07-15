package utils

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "0 seconds",
			duration: 0,
			expected: "0s",
		},
		{
			name:     "30 seconds",
			duration: 30 * time.Second,
			expected: "30s",
		},
		{
			name:     "59 seconds",
			duration: 59 * time.Second,
			expected: "59s",
		},
		{
			name:     "1 minute",
			duration: 1 * time.Minute,
			expected: "1m",
		},
		{
			name:     "5 minutes",
			duration: 5 * time.Minute,
			expected: "5m",
		},
		{
			name:     "59 minutes",
			duration: 59 * time.Minute,
			expected: "59m",
		},
		{
			name:     "1 hour",
			duration: 1 * time.Hour,
			expected: "1h",
		},
		{
			name:     "5 hours",
			duration: 5 * time.Hour,
			expected: "5h",
		},
		{
			name:     "23 hours",
			duration: 23 * time.Hour,
			expected: "23h",
		},
		{
			name:     "1 day",
			duration: 24 * time.Hour,
			expected: "1d",
		},
		{
			name:     "2 days",
			duration: 48 * time.Hour,
			expected: "2d",
		},
		{
			name:     "7 days",
			duration: 7 * 24 * time.Hour,
			expected: "7d",
		},
		{
			name:     "30 days",
			duration: 30 * 24 * time.Hour,
			expected: "30d",
		},
		{
			name:     "1 minute 30 seconds",
			duration: 1*time.Minute + 30*time.Second,
			expected: "1m",
		},
		{
			name:     "1 hour 30 minutes",
			duration: 1*time.Hour + 30*time.Minute,
			expected: "1h",
		},
		{
			name:     "1 day 5 hours",
			duration: 24*time.Hour + 5*time.Hour,
			expected: "1d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("FormatDuration(%v) = %s, want %s", tt.duration, result, tt.expected)
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "a greater than b",
			a:        10,
			b:        5,
			expected: 10,
		},
		{
			name:     "b greater than a",
			a:        3,
			b:        7,
			expected: 7,
		},
		{
			name:     "a equals b",
			a:        5,
			b:        5,
			expected: 5,
		},
		{
			name:     "negative numbers",
			a:        -3,
			b:        -7,
			expected: -3,
		},
		{
			name:     "mixed positive and negative",
			a:        -5,
			b:        3,
			expected: 3,
		},
		{
			name:     "zero values",
			a:        0,
			b:        0,
			expected: 0,
		},
		{
			name:     "one zero, one positive",
			a:        0,
			b:        5,
			expected: 5,
		},
		{
			name:     "one zero, one negative",
			a:        0,
			b:        -3,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Max(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Max(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "a less than b",
			a:        3,
			b:        7,
			expected: 3,
		},
		{
			name:     "b less than a",
			a:        10,
			b:        5,
			expected: 5,
		},
		{
			name:     "a equals b",
			a:        5,
			b:        5,
			expected: 5,
		},
		{
			name:     "negative numbers",
			a:        -3,
			b:        -7,
			expected: -7,
		},
		{
			name:     "mixed positive and negative",
			a:        -5,
			b:        3,
			expected: -5,
		},
		{
			name:     "zero values",
			a:        0,
			b:        0,
			expected: 0,
		},
		{
			name:     "one zero, one positive",
			a:        0,
			b:        5,
			expected: 0,
		},
		{
			name:     "one zero, one negative",
			a:        0,
			b:        -3,
			expected: -3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Min(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Min(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestFormatDurationEdgeCases(t *testing.T) {
	// Test edge cases for FormatDuration
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "exactly 1 minute",
			duration: 60 * time.Second,
			expected: "1m",
		},
		{
			name:     "exactly 1 hour",
			duration: 60 * time.Minute,
			expected: "1h",
		},
		{
			name:     "exactly 24 hours",
			duration: 24 * time.Hour,
			expected: "1d",
		},
		{
			name:     "very large duration",
			duration: 365 * 24 * time.Hour,
			expected: "365d",
		},
		{
			name:     "fractional seconds",
			duration: 500 * time.Millisecond,
			expected: "0s",
		},
		{
			name:     "negative duration",
			duration: -5 * time.Second,
			expected: "-5s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("FormatDuration(%v) = %s, want %s", tt.duration, result, tt.expected)
			}
		})
	}
}

func TestMaxMinWithLargeNumbers(t *testing.T) {
	// Test with very large numbers
	largeA := 2147483647  // Max int32
	largeB := 2147483646  // Max int32 - 1
	
	if Max(largeA, largeB) != largeA {
		t.Errorf("Max failed with large numbers")
	}
	
	if Min(largeA, largeB) != largeB {
		t.Errorf("Min failed with large numbers")
	}
	
	// Test with negative large numbers
	negLargeA := -2147483648  // Min int32
	negLargeB := -2147483647  // Min int32 + 1
	
	if Max(negLargeA, negLargeB) != negLargeB {
		t.Errorf("Max failed with negative large numbers")
	}
	
	if Min(negLargeA, negLargeB) != negLargeA {
		t.Errorf("Min failed with negative large numbers")
	}
}

func TestFormatDurationConsistency(t *testing.T) {
	// Test that formatting is consistent across similar durations
	durations := []time.Duration{
		1 * time.Second,
		2 * time.Second,
		59 * time.Second,
		60 * time.Second,
		61 * time.Second,
		119 * time.Second,
		120 * time.Second,
		121 * time.Second,
	}
	
	expected := []string{
		"1s",
		"2s", 
		"59s",
		"1m",
		"1m",
		"1m",
		"2m",
		"2m",
	}
	
	for i, duration := range durations {
		result := FormatDuration(duration)
		if result != expected[i] {
			t.Errorf("FormatDuration(%v) = %s, want %s", duration, result, expected[i])
		}
	}
}

func TestMaxMinSymmetry(t *testing.T) {
	// Test that Max and Min are symmetric
	testCases := []struct {
		a, b int
	}{
		{5, 10},
		{-5, -10},
		{0, 5},
		{-3, 7},
		{100, 100},
	}
	
	for _, tc := range testCases {
		// Max(a, b) should equal Max(b, a)
		if Max(tc.a, tc.b) != Max(tc.b, tc.a) {
			t.Errorf("Max is not symmetric: Max(%d, %d) != Max(%d, %d)", tc.a, tc.b, tc.b, tc.a)
		}
		
		// Min(a, b) should equal Min(b, a)
		if Min(tc.a, tc.b) != Min(tc.b, tc.a) {
			t.Errorf("Min is not symmetric: Min(%d, %d) != Min(%d, %d)", tc.a, tc.b, tc.b, tc.a)
		}
		
		// Max(a, b) should be >= Min(a, b)
		if Max(tc.a, tc.b) < Min(tc.a, tc.b) {
			t.Errorf("Max(%d, %d) < Min(%d, %d)", tc.a, tc.b, tc.a, tc.b)
		}
	}
}

func BenchmarkFormatDuration(b *testing.B) {
	durations := []time.Duration{
		30 * time.Second,
		5 * time.Minute,
		2 * time.Hour,
		3 * 24 * time.Hour,
	}
	
	for i := 0; i < b.N; i++ {
		for _, d := range durations {
			FormatDuration(d)
		}
	}
}

func BenchmarkMax(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Max(i, i+1)
	}
}

func BenchmarkMin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Min(i, i+1)
	}
}
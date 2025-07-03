package todo

import (
	"strings"
	"testing"
)

func TestValidateDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid description",
			description: "Valid description",
			wantErr:     false,
		},
		{
			name:        "empty description",
			description: "",
			wantErr:     true,
			errContains: "cannot be empty",
		},
		{
			name:        "whitespace only description",
			description: "   \t\n   ",
			wantErr:     true,
			errContains: "cannot be empty",
		},
		{
			name:        "description too long",
			description: strings.Repeat("a", 256),
			wantErr:     true,
			errContains: "too long",
		},
		{
			name:        "description at max length",
			description: strings.Repeat("a", 255),
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDescription(tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ValidateDescription() error = %v, want error containing %v", err, tt.errContains)
				}
			}
		})
	}
}

func TestParsePriority(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        Priority
		wantErr     bool
		errContains string
	}{
		{"low", "low", Low, false, ""},
		{"l", "l", Low, false, ""},
		{"1", "1", Low, false, ""},
		{"Low with caps", "Low", Low, false, ""},
		{"LOW all caps", "LOW", Low, false, ""},
		{"medium", "medium", Medium, false, ""},
		{"m", "m", Medium, false, ""},
		{"2", "2", Medium, false, ""},
		{"Medium with caps", "Medium", Medium, false, ""},
		{"MEDIUM all caps", "MEDIUM", Medium, false, ""},
		{"high", "high", High, false, ""},
		{"h", "h", High, false, ""},
		{"3", "3", High, false, ""},
		{"High with caps", "High", High, false, ""},
		{"HIGH all caps", "HIGH", High, false, ""},
		{"whitespace handling", "  medium  ", Medium, false, ""},
		{"invalid priority", "invalid", Medium, true, "must be one of"},
		{"empty string", "", Medium, true, "must be one of"},
		{"number out of range", "4", Medium, true, "must be one of"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePriority(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePriority() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ParsePriority() error = %v, want error containing %v", err, tt.errContains)
				}
			}
			if got != tt.want {
				t.Errorf("ParsePriority() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseStatus(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        Status
		wantErr     bool
		errContains string
	}{
		{"open", "open", Open, false, ""},
		{"o", "o", Open, false, ""},
		{"1", "1", Open, false, ""},
		{"Open with caps", "Open", Open, false, ""},
		{"OPEN all caps", "OPEN", Open, false, ""},
		{"done", "done", Done, false, ""},
		{"d", "d", Done, false, ""},
		{"2", "2", Done, false, ""},
		{"Done with caps", "Done", Done, false, ""},
		{"DONE all caps", "DONE", Done, false, ""},
		{"whitespace handling", "  done  ", Done, false, ""},
		{"invalid status", "invalid", Open, true, "must be one of"},
		{"empty string", "", Open, true, "must be one of"},
		{"number out of range", "3", Open, true, "must be one of"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStatus(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ParseStatus() error = %v, want error containing %v", err, tt.errContains)
				}
			}
			if got != tt.want {
				t.Errorf("ParseStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

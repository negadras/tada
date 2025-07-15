package components

import (
	"testing"

	"github.com/negadras/tada/internal/tui/styles"
	"github.com/negadras/tada/internal/tui/utils"
)

func TestForm_NewForm(t *testing.T) {
	styles := styles.DefaultStyles()
	keymap := utils.DefaultKeyMap()
	title := "Test Form"
	
	form := NewForm(styles, keymap, title)
	
	if form == nil {
		t.Fatal("NewForm returned nil")
	}
	
	if form.title != title {
		t.Errorf("Expected title '%s', got '%s'", title, form.title)
	}
	
	if form.styles != styles {
		t.Error("Form styles not set correctly")
	}
	
	if len(form.fields) != 0 {
		t.Errorf("Expected empty fields slice, got %d items", len(form.fields))
	}
	
	if form.focused != 0 {
		t.Errorf("Expected focused field to be 0, got %d", form.focused)
	}
	
	if form.submitted {
		t.Error("Expected form to not be submitted initially")
	}
	
	if form.cancelled {
		t.Error("Expected form to not be cancelled initially")
	}
}

func TestForm_AddField(t *testing.T) {
	form := NewForm(styles.DefaultStyles(), utils.DefaultKeyMap(), "Test Form")
	
	// Add required field
	form.AddField("Username", "Enter your username", true)
	
	if len(form.fields) != 1 {
		t.Errorf("Expected 1 field, got %d", len(form.fields))
	}
	
	field := form.fields[0]
	if field.Label != "Username" {
		t.Errorf("Expected label 'Username', got '%s'", field.Label)
	}
	
	if field.Placeholder != "Enter your username" {
		t.Errorf("Expected placeholder 'Enter your username', got '%s'", field.Placeholder)
	}
	
	if !field.Required {
		t.Error("Expected field to be required")
	}
	
	// Add optional field
	form.AddField("Email", "Enter your email", false)
	
	if len(form.fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(form.fields))
	}
	
	field2 := form.fields[1]
	if field2.Required {
		t.Error("Expected second field to be optional")
	}
}

func TestForm_GetValue(t *testing.T) {
	form := NewForm(styles.DefaultStyles(), utils.DefaultKeyMap(), "Test Form")
	
	// Add fields
	form.AddField("Field1", "Placeholder1", true)
	form.AddField("Field2", "Placeholder2", false)
	
	// Test getting values (should be empty initially)
	value1 := form.GetValue(0)
	if value1 != "" {
		t.Errorf("Expected empty value for field 0, got '%s'", value1)
	}
	
	value2 := form.GetValue(1)
	if value2 != "" {
		t.Errorf("Expected empty value for field 1, got '%s'", value2)
	}
	
	// Test getting value from invalid index
	valueInvalid := form.GetValue(5)
	if valueInvalid != "" {
		t.Errorf("Expected empty value for invalid index, got '%s'", valueInvalid)
	}
}

func TestForm_SetValue(t *testing.T) {
	form := NewForm(styles.DefaultStyles(), utils.DefaultKeyMap(), "Test Form")
	
	// Add fields
	form.AddField("Field1", "Placeholder1", true)
	form.AddField("Field2", "Placeholder2", false)
	
	// Set values
	form.SetValue(0, "Value 1")
	form.SetValue(1, "Value 2")
	
	// Test getting the set values
	value1 := form.GetValue(0)
	if value1 != "Value 1" {
		t.Errorf("Expected 'Value 1', got '%s'", value1)
	}
	
	value2 := form.GetValue(1)
	if value2 != "Value 2" {
		t.Errorf("Expected 'Value 2', got '%s'", value2)
	}
	
	// Test setting value at invalid index (should not panic)
	form.SetValue(5, "Invalid")
	
	// Test overwriting existing value
	form.SetValue(0, "New Value 1")
	value1Updated := form.GetValue(0)
	if value1Updated != "New Value 1" {
		t.Errorf("Expected 'New Value 1', got '%s'", value1Updated)
	}
}

func TestForm_Reset(t *testing.T) {
	form := NewForm(styles.DefaultStyles(), utils.DefaultKeyMap(), "Test Form")
	
	// Add fields and set values
	form.AddField("Field1", "Placeholder1", true)
	form.AddField("Field2", "Placeholder2", false)
	form.SetValue(0, "Value 1")
	form.SetValue(1, "Value 2")
	
	// Set form state
	form.submitted = true
	form.cancelled = true
	form.focused = 1
	
	// Reset form
	form.Reset()
	
	// Check that values are cleared
	value1 := form.GetValue(0)
	if value1 != "" {
		t.Errorf("Expected empty value after reset, got '%s'", value1)
	}
	
	value2 := form.GetValue(1)
	if value2 != "" {
		t.Errorf("Expected empty value after reset, got '%s'", value2)
	}
	
	// Check that state is reset
	if form.submitted {
		t.Error("Expected form to not be submitted after reset")
	}
	
	if form.cancelled {
		t.Error("Expected form to not be cancelled after reset")
	}
	
	if form.focused != 0 {
		t.Errorf("Expected focused field to be 0 after reset, got %d", form.focused)
	}
}

func TestForm_IsSubmitted(t *testing.T) {
	form := NewForm(styles.DefaultStyles(), utils.DefaultKeyMap(), "Test Form")
	
	// Initially not submitted
	if form.IsSubmitted() {
		t.Error("Expected form to not be submitted initially")
	}
	
	// Set as submitted
	form.submitted = true
	
	if !form.IsSubmitted() {
		t.Error("Expected form to be submitted")
	}
}

func TestForm_IsCancelled(t *testing.T) {
	form := NewForm(styles.DefaultStyles(), utils.DefaultKeyMap(), "Test Form")
	
	// Initially not cancelled
	if form.IsCancelled() {
		t.Error("Expected form to not be cancelled initially")
	}
	
	// Set as cancelled
	form.cancelled = true
	
	if !form.IsCancelled() {
		t.Error("Expected form to be cancelled")
	}
}

func TestForm_SetSize(t *testing.T) {
	form := NewForm(styles.DefaultStyles(), utils.DefaultKeyMap(), "Test Form")
	
	// Set size
	width := 80
	height := 24
	form.SetSize(width, height)
	
	if form.width != width {
		t.Errorf("Expected width %d, got %d", width, form.width)
	}
	
	if form.height != height {
		t.Errorf("Expected height %d, got %d", height, form.height)
	}
}

func TestForm_ValidationLogic(t *testing.T) {
	form := NewForm(styles.DefaultStyles(), utils.DefaultKeyMap(), "Test Form")
	
	// Add required and optional fields
	form.AddField("Required Field", "Enter required value", true)
	form.AddField("Optional Field", "Enter optional value", false)
	
	// Test validation logic (simulating form validation)
	tests := []struct {
		name               string
		requiredFieldValue string
		optionalFieldValue string
		expectedValid      bool
	}{
		{
			name:               "Valid form with both fields filled",
			requiredFieldValue: "Required Value",
			optionalFieldValue: "Optional Value",
			expectedValid:      true,
		},
		{
			name:               "Valid form with only required field filled",
			requiredFieldValue: "Required Value",
			optionalFieldValue: "",
			expectedValid:      true,
		},
		{
			name:               "Invalid form with required field empty",
			requiredFieldValue: "",
			optionalFieldValue: "Optional Value",
			expectedValid:      false,
		},
		{
			name:               "Invalid form with all fields empty",
			requiredFieldValue: "",
			optionalFieldValue: "",
			expectedValid:      false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form.Reset()
			form.SetValue(0, tt.requiredFieldValue)
			form.SetValue(1, tt.optionalFieldValue)
			
			// Simulate validation logic
			isValid := true
			for i, field := range form.fields {
				if field.Required && form.GetValue(i) == "" {
					isValid = false
					break
				}
			}
			
			if isValid != tt.expectedValid {
				t.Errorf("Expected validation result %v, got %v", tt.expectedValid, isValid)
			}
		})
	}
}

func TestForm_FocusNavigation(t *testing.T) {
	form := NewForm(styles.DefaultStyles(), utils.DefaultKeyMap(), "Test Form")
	
	// Add multiple fields
	form.AddField("Field1", "Placeholder1", true)
	form.AddField("Field2", "Placeholder2", false)
	form.AddField("Field3", "Placeholder3", true)
	
	// Test initial focus
	if form.focused != 0 {
		t.Errorf("Expected initial focus to be 0, got %d", form.focused)
	}
	
	// Test focus navigation logic (simulating Tab key behavior)
	nextFocus := (form.focused + 1) % len(form.fields)
	if nextFocus != 1 {
		t.Errorf("Expected next focus to be 1, got %d", nextFocus)
	}
	
	// Test focus at last field
	form.focused = len(form.fields) - 1
	nextFocus = (form.focused + 1) % len(form.fields)
	if nextFocus != 0 {
		t.Errorf("Expected focus to wrap to 0, got %d", nextFocus)
	}
	
	// Test previous focus navigation (simulating Shift+Tab)
	form.focused = 1
	prevFocus := (form.focused - 1 + len(form.fields)) % len(form.fields)
	if prevFocus != 0 {
		t.Errorf("Expected previous focus to be 0, got %d", prevFocus)
	}
	
	// Test previous focus from first field
	form.focused = 0
	prevFocus = (form.focused - 1 + len(form.fields)) % len(form.fields)
	if prevFocus != 2 {
		t.Errorf("Expected previous focus to wrap to 2, got %d", prevFocus)
	}
}

func TestForm_EmptyForm(t *testing.T) {
	form := NewForm(styles.DefaultStyles(), utils.DefaultKeyMap(), "Empty Form")
	
	// Test operations on empty form
	value := form.GetValue(0)
	if value != "" {
		t.Errorf("Expected empty value from empty form, got '%s'", value)
	}
	
	// Setting value on empty form should not panic
	form.SetValue(0, "Value")
	
	// Reset on empty form should not panic
	form.Reset()
	
	// Test that focus handling works with empty form
	if form.focused != 0 {
		t.Errorf("Expected focus to be 0 in empty form, got %d", form.focused)
	}
}
package validation

import "testing"

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email string
		want  bool
	}{
		{"test@example.com", true},
		{"user.name+tag@example.co.uk", true},
		{"invalid-email", false},
		{"@example.com", false},
		{"test@", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			if got := ValidateEmail(tt.email); got != tt.want {
				t.Errorf("ValidateEmail(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password string
		want     bool
	}{
		{"12345678", true},
		{"abcdefghij", true},
		{"1234567", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			if got := ValidatePassword(tt.password); got != tt.want {
				t.Errorf("ValidatePassword(%q) = %v, want %v", tt.password, got, tt.want)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		username string
		want     bool
	}{
		{"valid_username", true},
		{"user123", true},
		{"UsEr_Name", true},
		{"ab", false},                      // too short
		{"a_very_long_username_that_is_invalid", false}, // too long
		{"invalid-char", false},             // contains hyphen
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			if got := ValidateUsername(tt.username); got != tt.want {
				t.Errorf("ValidateUsername(%q) = %v, want %v", tt.username, got, tt.want)
			}
		})
	}
}

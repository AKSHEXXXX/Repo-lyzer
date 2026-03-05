package analyzer

import (
	"reflect"
	"testing"
)

func TestParseRequirementsTxt(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []Dependency
	}{
		{
			name:    "Simple package",
			content: "flask==2.0.0",
			expected: []Dependency{
				{Name: "flask", Version: "==2.0.0", Type: "production"},
			},
		},
		{
			name:    "Package with dot",
			content: "ruamel.yaml>=0.17.0",
			expected: []Dependency{
				{Name: "ruamel.yaml", Version: ">=0.17.0", Type: "production"},
			},
		},
		{
			name:    "Package with extras",
			content: "requests[security]==2.28.0",
			expected: []Dependency{
				{Name: "requests[security]", Version: "==2.28.0", Type: "production"},
			},
		},
		{
			name:    "Package with environment markers",
			content: "dataclasses; python_version < \"3.7\"",
			expected: []Dependency{
				{Name: "dataclasses", Version: "*", Type: "production"},
			},
		},
		{
			name:    "Package with version and marker",
			content: "requests==2.28.0 ; python_version > \"3.6\"",
			expected: []Dependency{
				{Name: "requests", Version: "==2.28.0", Type: "production"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := parseRequirementsTxt([]byte(tt.content))
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("parseRequirementsTxt() = %v, want %v", got, tt.expected)
			}
		})
	}
}

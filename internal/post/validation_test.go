package post

import (
	"strings"
	"testing"
)

func TestValidateContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "valid content",
			content: "Hello World",
			wantErr: false,
		},
		{
			name:    "empty content",
			content: "",
			wantErr: true,
		},
		{
			name:    "content too long",
			content: "a" + strings.Repeat("b", 300),
			wantErr: true,
		},
		{
			name:    "content with whitespace at start",
			content: "  Hello World",
			wantErr: false,
		},
		{
			name:    "content with only whitespace",
			content: "   ",
			wantErr: true,
		},
		{
			name:    "content with newline character",
			content: "Hello\nWorld",
			wantErr: false,
		},
		{
			name:    "content with special characters",
			content: "Hello!@#$%^&*()",
			wantErr: false,
		},
		{
			name:    "content with numbers",
			content: "Hello123",
			wantErr: false,
		},
		{
			name:    "single character content",
			content: "a",
			wantErr: false,
		},
		{
			name:    "content with 300 characters",
			content: strings.Repeat("a", 300),
			wantErr: false,
		},
		{
			name:    "content with 301 characters",
			content: strings.Repeat("a", 301),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateContent(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

package auth

import (
	"errors"
	"net/http"
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	tests := []struct {
		name        string
		headers     http.Header
		expectedKey string
		expectedErr error
	}{
		{
			name: "valid API key header",
			headers: http.Header{
				"Authorization": []string{"ApiKey my-secret-key-123"},
			},
			expectedKey: "my-secret-key-123",
			expectedErr: nil,
		},
		{
			name: "missing authorization header",
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			expectedKey: "",
			expectedErr: ErrNoAuthHeaderIncluded,
		},
		{
			name:        "empty headers",
			headers:     http.Header{},
			expectedKey: "",
			expectedErr: ErrNoAuthHeaderIncluded,
		},
		{
			name: "malformed header - only ApiKey prefix without key",
			headers: http.Header{
				"Authorization": []string{"ApiKey"},
			},
			expectedKey: "",
			expectedErr: errors.New("malformed authorization header"),
		},
		{
			name: "malformed header - wrong prefix",
			headers: http.Header{
				"Authorization": []string{"Bearer my-secret-key-123"},
			},
			expectedKey: "",
			expectedErr: errors.New("malformed authorization header"),
		},
		{
			name: "malformed header - no space after ApiKey",
			headers: http.Header{
				"Authorization": []string{"ApiKeimy-secret-key-123"},
			},
			expectedKey: "",
			expectedErr: errors.New("malformed authorization header"),
		},
		{
			name: "API key with multiple spaces - returns empty string",
			headers: http.Header{
				"Authorization": []string{"ApiKey    my-secret-key-123"},
			},
			expectedKey: "",
			expectedErr: nil,
		},
		{
			name: "API key with spaces - only first word returned",
			headers: http.Header{
				"Authorization": []string{"ApiKey key with spaces"},
			},
			expectedKey: "key",
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := GetAPIKey(tt.headers)

			if key != tt.expectedKey {
				t.Errorf("expected key %q, got %q", tt.expectedKey, key)
			}

			if tt.expectedErr == nil && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tt.expectedErr != nil && err == nil {
				t.Errorf("expected error %v, got nil", tt.expectedErr)
			}

			if tt.expectedErr != nil && err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

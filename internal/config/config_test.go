package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		wantC   *Config
		name    string
		wantErr bool
	}{
		{
			name:    "OK",
			wantC:   &Config{Host: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: "db.json", Secret: "fortytwo"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := New()
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantC, gotC)
		})
	}
}

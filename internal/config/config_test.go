package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		wantC   Сonfig
		wantErr bool
	}{
		{
			name:    "OK",
			wantC:   Сonfig{Host: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: "db.json"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := Parse()
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantC, gotC)
		})
	}
}

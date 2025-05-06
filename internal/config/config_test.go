package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockConfig struct {
	Host            string `json:"Host"`
	BaseURL         string `json:"BaseURL"`
	FileStoragePath string `json:"FileStoragePath"`
	DatabaseDSN     string `json:"DatabaseDSN"`
	Secret          string `json:"Secret"`
	TLSCertPath     string `json:"TLSCertPath"`
	TLSKeyPath      string `json:"TLSKeyPath"`
	TrustedSubnet   string `json:"TrustedSubnet"`
	EnableHTTPS     bool   `json:"EnableHTTPS"`
}

func TestLoadJSONConfig(t *testing.T) {
	// Временный JSON-файл
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.json")
	badfile := "notExistFile.json"

	jsonContent := `{
		"EnableHTTPS": true,
		"TLSCertPath": "certs/cert.crt",
		"TLSKeyPath": "certs/cert.key",
		"Host": ":8080",
		"BaseURL": "https://localhost:8080",
		"FileStoragePath": "./tmp/shorten-db.json",
		"DatabaseDSN": "",
		"Secret": "super-secret",
		"TrustedSubnet": "127.0.0.0/24"
	}`

	badJSONContent := `{ret": "super-secretbsdb2 }`

	var cfg mockConfig

	err := LoadJSONConfig(badfile, &cfg)
	require.Error(t, err)

	err = os.WriteFile(tmpFile, []byte(badJSONContent), 0600)
	require.NoError(t, err)

	err = LoadJSONConfig(tmpFile, &cfg)
	require.Error(t, err)

	err = os.WriteFile(tmpFile, []byte(jsonContent), 0600)
	require.NoError(t, err)

	err = LoadJSONConfig(tmpFile, &cfg)
	require.NoError(t, err)

	assert.Equal(t, true, cfg.EnableHTTPS)
	assert.Equal(t, "certs/cert.crt", cfg.TLSCertPath)
	assert.Equal(t, "certs/cert.key", cfg.TLSKeyPath)
	assert.Equal(t, ":8080", cfg.Host)
	assert.Equal(t, "https://localhost:8080", cfg.BaseURL)
	assert.Equal(t, "./tmp/shorten-db.json", cfg.FileStoragePath)
	assert.Equal(t, "", cfg.DatabaseDSN)
	assert.Equal(t, "super-secret", cfg.Secret)
	assert.Equal(t, "127.0.0.0/24", cfg.TrustedSubnet)
}

func TestParse(t *testing.T) {
	tests := []struct {
		wantC   *Config
		name    string
		wantErr bool
	}{
		{
			name:    "OK",
			wantC:   &Config{EnableHTTPS: false, TLSCertPath: "certs/cert.crt", TLSKeyPath: "certs/cert.key", Config: "", Host: "localhost:8080", BaseURL: "http://localhost:8080", FileStoragePath: "db.json", DatabaseDSN: "", Secret: "fortytwo", TrustedSubnet: "127.0.0.0/24"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := New()
			t.Log(gotC)
			t.Log(tt.wantC)
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantC, gotC)
		})
	}
}

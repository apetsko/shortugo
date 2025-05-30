package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func random(max int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		log.Println(err)
	}
	return nBig.Int64()
}

func generateURLS(count int) []string {
	urls := make([]string, count)
	protocols := []string{
		"http://",
		"https://",
		// "ftp://",
		// "ftps://",
		// "sftp://",
		// "mailto://",
		// "telnet://",
		// "file://",
		// "data://",
		// "ws://",
		// "wss://",
		// "bluetooth://",
	}

	for i := range urls {
		u := fmt.Sprintf("%s%s.%s/%s/%s",
			protocols[random(int64(len(protocols)))], // protocol
			randomString(random(int64(10))+1),        // 2domain
			randomString(random(int64(3))+2),         // 1domain
			randomString(random(int64(9))+2),         // path1
			randomString(random(int64(13))+4),        // path2
		)
		urls[i] = u
	}

	return urls
}

func randomString(length int64) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[random(int64(len(charset)))]
	}
	return string(b)
}

func Test_Generate(t *testing.T) {
	m := make(map[string]string)
	urls := generateURLS(10)
	for i, u := range urls {
		t.Run(fmt.Sprintf("URL #%d", i), func(t *testing.T) {
			IDlen := 8
			ID := GenerateID(u, IDlen)

			ur, ok := m[ID]
			require.Equal(t, false, ok)
			require.Equal(t, false, ur == u)

			m[ID] = u
		})
	}
}
func TestGenerateID(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		length   int
	}{
		{
			name:     "Generate 8-char ID",
			input:    "test_input",
			length:   8,
			expected: "lSgi3mpi",
		},
		{
			name:     "Generate 12-char ID",
			input:    "another_input",
			length:   12,
			expected: "2OUOxCc_9wx6",
		},
		{
			name:     "Generate 6-char ID",
			input:    "short",
			length:   6,
			expected: "-bAHi1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id := GenerateID(tc.input, tc.length)
			assert.Len(t, id, tc.length)
			assert.Equal(t, id, tc.expected)
		})
	}
}

func TestGenerateUserID(t *testing.T) {
	testCases := []struct {
		name   string
		length int
	}{
		{
			name:   "Generate 16-byte User ID",
			length: 16,
		},
		{
			name:   "Generate 32-byte User ID",
			length: 32,
		},
		{
			name:   "Generate 8-byte User ID",
			length: 8,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := GenerateUserID(tc.length)
			assert.NoError(t, err)
			assert.Len(t, id, tc.length*2) // hex encoding удваивает длину
		})
	}
}

func TestValidateStruct(t *testing.T) {
	type TestStruct struct {
		Name string `validate:"required"`
		Age  int    `validate:"min=18"`
	}

	testCases := []struct {
		name      string
		input     TestStruct
		errFields []string
		expectErr bool
	}{
		{
			name:      "Invalid: empty Name",
			input:     TestStruct{Name: "", Age: 20},
			expectErr: true,
			errFields: []string{"Name"},
		},
		{
			name:      "Invalid: Age < 18",
			input:     TestStruct{Name: "Alice", Age: 16},
			expectErr: true,
			errFields: []string{"Age"},
		},
		{
			name:      "Invalid: empty Name and Age < 18",
			input:     TestStruct{Name: "", Age: 16},
			expectErr: true,
			errFields: []string{"Name", "Age"},
		},
		{
			name:      "Valid struct",
			input:     TestStruct{Name: "Bob", Age: 25},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateStruct(tc.input)

			if tc.expectErr {
				assert.Error(t, err)
				for _, field := range tc.errFields {
					assert.Contains(t, err.Error(), field)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

package utils

import (
	"errors"
	"fmt"
	"testing"
)

func Test_Generate(t *testing.T) {
	t.Run("simple Tesst", func(t *testing.T) {
		u := "asdfasdfsdf"
		if ID, err := Generate(u); ID != "EwHXdJfB" || err != nil {
			t.Errorf("%q, %q := Generate(%s)", ID, err.Error(), u)
		}
	})
}

func Test_FullURL(t *testing.T) {
	type want struct {
		URL     string
		err     error
		wantErr bool
	}

	tests := []struct {
		want    want
		baseURL string
		id      string
	}{
		{want: want{"", errors.New("empty id"), true}, id: "", baseURL: "localhost:8080"},
		{want: want{"http://localhost:8080/EwHXdJfB", nil, false}, id: "EwHXdJfB"},
	}

	for i, tt := range tests {
		t.Run((fmt.Sprintf("%v", i)), func(t *testing.T) {
			if URL, err := FullUrl(tt.baseURL, tt.id); (err != nil) != tt.want.wantErr {
				t.Errorf("%q, %q := FullUrl(%s)", URL, err.Error(), tt.id)
			}
		})
	}
}

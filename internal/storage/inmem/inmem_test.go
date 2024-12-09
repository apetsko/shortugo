package inmem

import (
	"testing"

	"github.com/apetsko/shortugo/internal/utils"
)

func Test_Put(t *testing.T) {
	im := New()

	test := struct {
		URL     string
		wantID  string
		wantErr bool
	}{
		URL:     "https://practicum.yandex.ru/",
		wantErr: false,
	}
	t.Run("test_put", func(t *testing.T) {

		ID := utils.Generate(test.URL)

		im.data[ID] = test.URL

		if v, ok := im.data[ID]; !ok || test.URL != v {
			t.Errorf("%q, %t := im.data[%q]", v, ok, test.wantID)
		}
	})
}

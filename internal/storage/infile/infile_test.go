package infile

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Put(t *testing.T) {
	ifile, err := New("db_test.json")
	require.NoError(t, err)

	tests := []struct {
		URL string
		ID  string
	}{
		{"mailto://EBlI.LUcE/nGW/CnKgralWM", "EVvMeswX"},
		{"data://bNZlqPkX.zPr/AOYjayx/RXDZywCjbH", "zrWsrYVK"},
		{"ftps://QhPSk.SERo/ASOuRTdh/XuXCUVcR", "WrBTersI"},
		{"http://hwr.DqhY/qRpylA/BrBUqXwraQX", "IZBF3Drj"},
		{"file://rSX.gQs/AoJCRUFJbS/HbkVkdDhHkSakU", "B-_ig72W"},
		{"file://c.Hh/Oo/cAWXXgykO", "ih4UOFRN"},
		{"http://rfcv.yZ/djwBnRy/GRvWfxKARJXqiIS", "CnhlRf81"},
		{"sftp://zvJXD.xR/lUTNLwCMuL/ACaRzHI", "oSyiotBD"},
		{"ws://SAZCfOUSn.qxaU/tj/TIdK", "7la40tTW"},
		{"file://IyZL.go/YfaSpOpqhN/XfWd", "7HVUuC38"},
		{"telnet://npLzsEwn.KTR/XLv/gYhEqqdTTCUdpEjE", "_QDwIZ8V"},
		{"ftps://PlqcUsANz.fn/wpSOrY/NVHIDGTbCVUSL", "JJd8nofa"},
		{"file://WLCHVIgAk.Nc/gAqCVuw/GBZaquHPx", "SVKhwBjn"},
		{"bluetooth://qtuD.eT/OugB/XeohyIVkj", "jzLEbSpd"},
		{"file://hya.jrqF/smmqgM/GJeaDJOYx", "UrqyUbm_"},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test_put #%d", i), func(t *testing.T) {
			err := ifile.Put(test.ID, test.URL)
			require.NoError(t, err)

			v, err := ifile.Get(test.ID)
			require.NoError(t, err)
			assert.Equal(t, v, test.URL)
		})
	}
}

func Test_Get(t *testing.T) {
	ifile, err := New("db_test.json")
	require.NoError(t, err)

	tests := []struct {
		URL string
		ID  string
	}{
		{"mailto://EBlI.LUcE/nGW/CnKgralWM", "EVvMeswX"},
		{"data://bNZlqPkX.zPr/AOYjayx/RXDZywCjbH", "zrWsrYVK"},
		{"ftps://QhPSk.SERo/ASOuRTdh/XuXCUVcR", "WrBTersI"},
		{"http://hwr.DqhY/qRpylA/BrBUqXwraQX", "IZBF3Drj"},
		{"file://rSX.gQs/AoJCRUFJbS/HbkVkdDhHkSakU", "B-_ig72W"},
		{"file://c.Hh/Oo/cAWXXgykO", "ih4UOFRN"},
		{"http://rfcv.yZ/djwBnRy/GRvWfxKARJXqiIS", "CnhlRf81"},
		{"sftp://zvJXD.xR/lUTNLwCMuL/ACaRzHI", "oSyiotBD"},
		{"ws://SAZCfOUSn.qxaU/tj/TIdK", "7la40tTW"},
		{"file://IyZL.go/YfaSpOpqhN/XfWd", "7HVUuC38"},
		{"telnet://npLzsEwn.KTR/XLv/gYhEqqdTTCUdpEjE", "_QDwIZ8V"},
		{"ftps://PlqcUsANz.fn/wpSOrY/NVHIDGTbCVUSL", "JJd8nofa"},
		{"file://WLCHVIgAk.Nc/gAqCVuw/GBZaquHPx", "SVKhwBjn"},
		{"bluetooth://qtuD.eT/OugB/XeohyIVkj", "jzLEbSpd"},
		{"file://hya.jrqF/smmqgM/GJeaDJOYx", "UrqyUbm_"},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test_put #%d", i), func(t *testing.T) {
			u, err := ifile.Get(test.ID)
			require.NoError(t, err)
			assert.Equal(t, u, test.URL)
		})
	}
}

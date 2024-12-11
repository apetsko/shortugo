package inmem

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Put(t *testing.T) {
	im := New()

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
			err := im.Put(test.ID, test.URL)
			require.NoError(t, err)

			v, ok := im.data[test.ID]
			require.Equal(t, ok, true)
			assert.Equal(t, v, test.URL)
		})
	}
}

func Test_Get(t *testing.T) {
	im := New()
	im.data = map[string]string{
		"EVvMeswX": "mailto://EBlI.LUcE/nGW/CnKgralWM",
		"zrWsrYVK": "data://bNZlqPkX.zPr/AOYjayx/RXDZywCjbH",
		"WrBTersI": "ftps://QhPSk.SERo/ASOuRTdh/XuXCUVcR",
		"IZBF3Drj": "http://hwr.DqhY/qRpylA/BrBUqXwraQX",
		"B-_ig72W": "file://rSX.gQs/AoJCRUFJbS/HbkVkdDhHkSakU",
		"ih4UOFRN": "file://c.Hh/Oo/cAWXXgykO",
		"CnhlRf81": "http://rfcv.yZ/djwBnRy/GRvWfxKARJXqiIS",
		"oSyiotBD": "sftp://zvJXD.xR/lUTNLwCMuL/ACaRzHI",
		"7la40tTW": "ws://SAZCfOUSn.qxaU/tj/TIdK",
		"7HVUuC38": "file://IyZL.go/YfaSpOpqhN/XfWd",
		"_QDwIZ8V": "telnet://npLzsEwn.KTR/XLv/gYhEqqdTTCUdpEjE",
		"JJd8nofa": "ftps://PlqcUsANz.fn/wpSOrY/NVHIDGTbCVUSL",
		"SVKhwBjn": "file://WLCHVIgAk.Nc/gAqCVuw/GBZaquHPx",
		"jzLEbSpd": "bluetooth://qtuD.eT/OugB/XeohyIVkj",
		"UrqyUbm_": "file://hya.jrqF/smmqgM/GJeaDJOYx",
	}

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

			u, err := im.Get(test.ID)
			require.NoError(t, err)
			assert.Equal(t, u, test.URL)
		})
	}
}

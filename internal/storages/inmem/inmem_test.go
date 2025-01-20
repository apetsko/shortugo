package inmem

import (
	"context"
	"fmt"
	"testing"

	"github.com/apetsko/shortugo/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Put(t *testing.T) {
	im := New()

	tests := []models.URLRecord{
		{URL: "mailto://EBlI.LUcE/nGW/CnKgralWM", ID: "EVvMeswX"},
		{URL: "data://bNZlqPkX.zPr/AOYjayx/RXDZywCjbH", ID: "zrWsrYVK"},
		{URL: "ftps://QhPSk.SERo/ASOuRTdh/XuXCUVcR", ID: "WrBTersI"},
		{URL: "http://hwr.DqhY/qRpylA/BrBUqXwraQX", ID: "IZBF3Drj"},
		{URL: "file://rSX.gQs/AoJCRUFJbS/HbkVkdDhHkSakU", ID: "B -_ig72W"},
		{URL: "file://c.Hh/Oo/cAWXXgykO", ID: "ih4UOFRN"},
		{URL: "http://rfcv.yZ/djwBnRy/GRvWfxKARJXqiIS", ID: "CnhlRf81"},
		{URL: "sftp://zvJXD.xR/lUTNLwCMuL/ACaRzHI", ID: "oSyiotBD"},
		{URL: "ws://SAZCfOUSn.qxaU/tj/TIdK", ID: "7la40tTW"},
		{URL: "file://IyZL.go/YfaSpOpqhN/XfWd", ID: "7HVUuC38"},
		{URL: "telnet://npLzsEwn.KTR/XLv/gYhEqqdTTCUdpEjE", ID: "_QDwIZ8V"},
		{URL: "ftps://PlqcUsANz.fn/wpSOrY/NVHIDGTbCVUSL", ID: "JJd8nofa"},
		{URL: "file://WLCHVIgAk.Nc/gAqCVuw/GBZaquHPx", ID: "SVKhwBjn"},
		{URL: "bluetooth://qtuD.eT/OugB/XeohyIVkj", ID: "jzLEbSpd"},
		{URL: "file://hya.jrqF/smmqgM/GJeaDJOYx", ID: "UrqyUbm_"},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test_put #%d", i), func(t *testing.T) {
			err := im.Put(context.Background(), test)
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

	tests := []models.URLRecord{
		{URL: "mailto://EBlI.LUcE/nGW/CnKgralWM", ID: "EVvMeswX"},
		{URL: "data://bNZlqPkX.zPr/AOYjayx/RXDZywCjbH", ID: "zrWsrYVK"},
		{URL: "ftps://QhPSk.SERo/ASOuRTdh/XuXCUVcR", ID: "WrBTersI"},
		{URL: "http://hwr.DqhY/qRpylA/BrBUqXwraQX", ID: "IZBF3Drj"},
		{URL: "file://rSX.gQs/AoJCRUFJbS/HbkVkdDhHkSakU", ID: "B-_ig72W"},
		{URL: "file://c.Hh/Oo/cAWXXgykO", ID: "ih4UOFRN"},
		{URL: "http://rfcv.yZ/djwBnRy/GRvWfxKARJXqiIS", ID: "CnhlRf81"},
		{URL: "sftp://zvJXD.xR/lUTNLwCMuL/ACaRzHI", ID: "oSyiotBD"},
		{URL: "ws://SAZCfOUSn.qxaU/tj/TIdK", ID: "7la40tTW"},
		{URL: "file://IyZL.go/YfaSpOpqhN/XfWd", ID: "7HVUuC38"},
		{URL: "telnet://npLzsEwn.KTR/XLv/gYhEqqdTTCUdpEjE", ID: "_QDwIZ8V"},
		{URL: "ftps://PlqcUsANz.fn/wpSOrY/NVHIDGTbCVUSL", ID: "JJd8nofa"},
		{URL: "file://WLCHVIgAk.Nc/gAqCVuw/GBZaquHPx", ID: "SVKhwBjn"},
		{URL: "bluetooth://qtuD.eT/OugB/XeohyIVkj", ID: "jzLEbSpd"},
		{URL: "file://hya.jrqF/smmqgM/GJeaDJOYx", ID: "UrqyUbm_"},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test_put #%d", i), func(t *testing.T) {
			ctx := context.Background()
			u, err := im.Get(ctx, test.ID)
			require.NoError(t, err)
			assert.Equal(t, u, test.URL)
		})
	}
}

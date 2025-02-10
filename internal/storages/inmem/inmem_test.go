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
		{UserID: "22", URL: "mailto://EBlI.LUcE/nGW/CnKgralWM", ID: "EVvMeswX"},
		{UserID: "22", URL: "data://bNZlqPkX.zPr/AOYjayx/RXDZywCjbH", ID: "zrWsrYVK"},
		{UserID: "22", URL: "ftps://QhPSk.SERo/ASOuRTdh/XuXCUVcR", ID: "WrBTersI"},
		{UserID: "22", URL: "http://hwr.DqhY/qRpylA/BrBUqXwraQX", ID: "IZBF3Drj"},
		{UserID: "22", URL: "file://rSX.gQs/AoJCRUFJbS/HbkVkdDhHkSakU", ID: "B-_ig72W"},
		{UserID: "22", URL: "file://c.Hh/Oo/cAWXXgykO", ID: "ih4UOFRN"},
		{UserID: "22", URL: "http://rfcv.yZ/djwBnRy/GRvWfxKARJXqiIS", ID: "CnhlRf81"},
		{UserID: "22", URL: "sftp://zvJXD.xR/lUTNLwCMuL/ACaRzHI", ID: "oSyiotBD"},
		{UserID: "1", URL: "ws://SAZCfOUSn.qxaU/tj/TIdK", ID: "7la40tTW"},
		{UserID: "1", URL: "file://IyZL.go/YfaSpOpqhN/XfWd", ID: "7HVUuC38"},
		{UserID: "1", URL: "telnet://npLzsEwn.KTR/XLv/gYhEqqdTTCUdpEjE", ID: "_QDwIZ8V"},
		{UserID: "1", URL: "ftps://PlqcUsANz.fn/wpSOrY/NVHIDGTbCVUSL", ID: "JJd8nofa"},
		{UserID: "1", URL: "file://WLCHVIgAk.Nc/gAqCVuw/GBZaquHPx", ID: "SVKhwBjn"},
		{UserID: "1", URL: "bluetooth://qtuD.eT/OugB/XeohyIVkj", ID: "jzLEbSpd"},
		{UserID: "1", URL: "file://hya.jrqF/smmqgM/GJeaDJOYx", ID: "UrqyUbm_"},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test_put #%d", i), func(t *testing.T) {
			err := im.Put(context.Background(), test)
			require.NoError(t, err)

			v, ok := im.byID[test.ID]
			require.Equal(t, ok, true)
			assert.Equal(t, v, test)
		})
	}
}

func Test_Get(t *testing.T) {
	im := New()
	im.byID = map[string]models.URLRecord{
		"EVvMeswX": {UserID: "22", URL: "mailto://EBlI.LUcE/nGW/CnKgralWM", ID: "EVvMeswX"},
		"zrWsrYVK": {UserID: "22", URL: "data://bNZlqPkX.zPr/AOYjayx/RXDZywCjbH", ID: "zrWsrYVK"},
		"WrBTersI": {UserID: "22", URL: "ftps://QhPSk.SERo/ASOuRTdh/XuXCUVcR", ID: "WrBTersI"},
		"IZBF3Drj": {UserID: "22", URL: "http://hwr.DqhY/qRpylA/BrBUqXwraQX", ID: "IZBF3Drj"},
		"B-_ig72W": {UserID: "22", URL: "file://rSX.gQs/AoJCRUFJbS/HbkVkdDhHkSakU", ID: "B-_ig72W"},
		"ih4UOFRN": {UserID: "22", URL: "file://c.Hh/Oo/cAWXXgykO", ID: "ih4UOFRN"},
		"CnhlRf81": {UserID: "22", URL: "http://rfcv.yZ/djwBnRy/GRvWfxKARJXqiIS", ID: "CnhlRf81"},
		"oSyiotBD": {UserID: "22", URL: "sftp://zvJXD.xR/lUTNLwCMuL/ACaRzHI", ID: "oSyiotBD"},
		"7la40tTW": {UserID: "1", URL: "ws://SAZCfOUSn.qxaU/tj/TIdK", ID: "7la40tTW"},
		"7HVUuC38": {UserID: "1", URL: "file://IyZL.go/YfaSpOpqhN/XfWd", ID: "7HVUuC38"},
		"_QDwIZ8V": {UserID: "1", URL: "telnet://npLzsEwn.KTR/XLv/gYhEqqdTTCUdpEjE", ID: "_QDwIZ8V"},
		"JJd8nofa": {UserID: "1", URL: "ftps://PlqcUsANz.fn/wpSOrY/NVHIDGTbCVUSL", ID: "JJd8nofa"},
		"SVKhwBjn": {UserID: "1", URL: "file://WLCHVIgAk.Nc/gAqCVuw/GBZaquHPx", ID: "SVKhwBjn"},
		"jzLEbSpd": {UserID: "1", URL: "bluetooth://qtuD.eT/OugB/XeohyIVkj", ID: "jzLEbSpd"},
		"UrqyUbm_": {UserID: "1", URL: "file://hya.jrqF/smmqgM/GJeaDJOYx", ID: "UrqyUbm_"},
	}

	tests := []models.URLRecord{
		{UserID: "22", URL: "mailto://EBlI.LUcE/nGW/CnKgralWM", ID: "EVvMeswX"},
		{UserID: "22", URL: "data://bNZlqPkX.zPr/AOYjayx/RXDZywCjbH", ID: "zrWsrYVK"},
		{UserID: "22", URL: "ftps://QhPSk.SERo/ASOuRTdh/XuXCUVcR", ID: "WrBTersI"},
		{UserID: "22", URL: "http://hwr.DqhY/qRpylA/BrBUqXwraQX", ID: "IZBF3Drj"},
		{UserID: "22", URL: "file://rSX.gQs/AoJCRUFJbS/HbkVkdDhHkSakU", ID: "B-_ig72W"},
		{UserID: "22", URL: "file://c.Hh/Oo/cAWXXgykO", ID: "ih4UOFRN"},
		{UserID: "22", URL: "http://rfcv.yZ/djwBnRy/GRvWfxKARJXqiIS", ID: "CnhlRf81"},
		{UserID: "22", URL: "sftp://zvJXD.xR/lUTNLwCMuL/ACaRzHI", ID: "oSyiotBD"},
		{UserID: "1", URL: "ws://SAZCfOUSn.qxaU/tj/TIdK", ID: "7la40tTW"},
		{UserID: "1", URL: "file://IyZL.go/YfaSpOpqhN/XfWd", ID: "7HVUuC38"},
		{UserID: "1", URL: "telnet://npLzsEwn.KTR/XLv/gYhEqqdTTCUdpEjE", ID: "_QDwIZ8V"},
		{UserID: "1", URL: "ftps://PlqcUsANz.fn/wpSOrY/NVHIDGTbCVUSL", ID: "JJd8nofa"},
		{UserID: "1", URL: "file://WLCHVIgAk.Nc/gAqCVuw/GBZaquHPx", ID: "SVKhwBjn"},
		{UserID: "1", URL: "bluetooth://qtuD.eT/OugB/XeohyIVkj", ID: "jzLEbSpd"},
		{UserID: "1", URL: "file://hya.jrqF/smmqgM/GJeaDJOYx", ID: "UrqyUbm_"},
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

func TestStorage_GetAllLinksByUser000ID(t *testing.T) {
	im := New()
	im.byUserID = map[string][]models.URLRecord{
		"22": {
			{UserID: "22", URL: "mailto://EBlI.LUcE/nGW/CnKgralWM", ID: "EVvMeswX"},
			{UserID: "22", URL: "data://bNZlqPkX.zPr/AOYjayx/RXDZywCjbH", ID: "zrWsrYVK"},
			{UserID: "22", URL: "ftps://QhPSk.SERo/ASOuRTdh/XuXCUVcR", ID: "WrBTersI"},
			{UserID: "22", URL: "http://hwr.DqhY/qRpylA/BrBUqXwraQX", ID: "IZBF3Drj"},
			{UserID: "22", URL: "file://rSX.gQs/AoJCRUFJbS/HbkVkdDhHkSakU", ID: "B-_ig72W"},
			{UserID: "22", URL: "file://c.Hh/Oo/cAWXXgykO", ID: "ih4UOFRN"},
			{UserID: "22", URL: "http://rfcv.yZ/djwBnRy/GRvWfxKARJXqiIS", ID: "CnhlRf81"},
			{UserID: "22", URL: "sftp://zvJXD.xR/lUTNLwCMuL/ACaRzHI", ID: "oSyiotBD"},
		},
		"1": {
			{UserID: "1", URL: "ws://SAZCfOUSn.qxaU/tj/TIdK", ID: "7la40tTW"},
			{UserID: "1", URL: "file://IyZL.go/YfaSpOpqhN/XfWd", ID: "7HVUuC38"},
			{UserID: "1", URL: "telnet://npLzsEwn.KTR/XLv/gYhEqqdTTCUdpEjE", ID: "_QDwIZ8V"},
			{UserID: "1", URL: "ftps://PlqcUsANz.fn/wpSOrY/NVHIDGTbCVUSL", ID: "JJd8nofa"},
			{UserID: "1", URL: "file://WLCHVIgAk.Nc/gAqCVuw/GBZaquHPx", ID: "SVKhwBjn"},
			{UserID: "1", URL: "bluetooth://qtuD.eT/OugB/XeohyIVkj", ID: "jzLEbSpd"},
			{UserID: "1", URL: "file://hya.jrqF/smmqgM/GJeaDJOYx", ID: "UrqyUbm_"},
		}}

	tests := map[string][]models.URLRecord{
		"22": {
			{UserID: "22", URL: "mailto://EBlI.LUcE/nGW/CnKgralWM", ID: "/EVvMeswX"},
			{UserID: "22", URL: "data://bNZlqPkX.zPr/AOYjayx/RXDZywCjbH", ID: "/zrWsrYVK"},
			{UserID: "22", URL: "ftps://QhPSk.SERo/ASOuRTdh/XuXCUVcR", ID: "/WrBTersI"},
			{UserID: "22", URL: "http://hwr.DqhY/qRpylA/BrBUqXwraQX", ID: "/IZBF3Drj"},
			{UserID: "22", URL: "file://rSX.gQs/AoJCRUFJbS/HbkVkdDhHkSakU", ID: "/B-_ig72W"},
			{UserID: "22", URL: "file://c.Hh/Oo/cAWXXgykO", ID: "/ih4UOFRN"},
			{UserID: "22", URL: "http://rfcv.yZ/djwBnRy/GRvWfxKARJXqiIS", ID: "/CnhlRf81"},
			{UserID: "22", URL: "sftp://zvJXD.xR/lUTNLwCMuL/ACaRzHI", ID: "/oSyiotBD"},
		},
		"1": {
			{UserID: "1", URL: "ws://SAZCfOUSn.qxaU/tj/TIdK", ID: "/7la40tTW"},
			{UserID: "1", URL: "file://IyZL.go/YfaSpOpqhN/XfWd", ID: "/7HVUuC38"},
			{UserID: "1", URL: "telnet://npLzsEwn.KTR/XLv/gYhEqqdTTCUdpEjE", ID: "/_QDwIZ8V"},
			{UserID: "1", URL: "ftps://PlqcUsANz.fn/wpSOrY/NVHIDGTbCVUSL", ID: "/JJd8nofa"},
			{UserID: "1", URL: "file://WLCHVIgAk.Nc/gAqCVuw/GBZaquHPx", ID: "/SVKhwBjn"},
			{UserID: "1", URL: "bluetooth://qtuD.eT/OugB/XeohyIVkj", ID: "/jzLEbSpd"},
			{UserID: "1", URL: "file://hya.jrqF/smmqgM/GJeaDJOYx", ID: "/UrqyUbm_"},
		}}

	for userID := range tests {
		t.Run(userID, func(t *testing.T) {
			ctx := context.Background()
			gotRr, err := im.ListLinksByUserID(ctx, "", userID)
			require.NoError(t, err)
			assert.Equalf(t, tests[userID], gotRr, "ListLinksByUserID(%v, %v)", ctx, userID)
		})
	}
}

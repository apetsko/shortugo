package inmem

import (
	"context"
	"fmt"
	"testing"

	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/shared"
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

func TestStorage_PutAndGet(t *testing.T) {
	store := New()
	ctx := context.Background()

	testCases := []struct {
		name    string
		wantURL string
		wantErr error
		record  models.URLRecord
	}{
		{
			name: "Valid record",
			record: models.URLRecord{
				ID:     "short1",
				URL:    "https://example.com",
				UserID: "user1",
			},
			wantURL: "https://example.com",
			wantErr: nil,
		},
		{
			name: "Another valid record",
			record: models.URLRecord{
				ID:     "short2",
				URL:    "https://test.com",
				UserID: "user2",
			},
			wantURL: "https://test.com",
			wantErr: nil,
		},
		{
			name:    "Not found",
			record:  models.URLRecord{},
			wantURL: "",
			wantErr: shared.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.record.ID != "" {
				require.NoError(t, store.Put(ctx, tc.record))
			}

			gotURL, err := store.Get(ctx, tc.record.ID)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantURL, gotURL)
			}
		})
	}
}

func TestStorage_DeleteUserURLs(t *testing.T) {
	store := New()
	ctx := context.Background()

	records := []models.URLRecord{
		{ID: "short1", URL: "https://example.com", UserID: "user1"},
		{ID: "short2", URL: "https://test.com", UserID: "user1"},
	}

	for _, r := range records {
		require.NoError(t, store.Put(ctx, r))
	}

	testCases := []struct {
		wantErr error
		name    string
		userID  string
		ids     []string
	}{
		{
			name:   "Delete existing URL",
			ids:    []string{"short1"},
			userID: "user1",
		},
		{
			name:    "Delete non-existing URL",
			ids:     []string{"short3"},
			userID:  "user1",
			wantErr: shared.ErrNotFound,
		},
		{
			name:   "Delete with wrong userID",
			ids:    []string{"short2"},
			userID: "user2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := store.DeleteUserURLs(ctx, tc.ids, tc.userID)
			assert.NoError(t, err)

			for _, id := range tc.ids {
				_, err := store.Get(ctx, id)
				if tc.userID == "user1" && id == "short1" {
					assert.ErrorIs(t, err, shared.ErrGone)
				} else if id == "short3" {
					assert.ErrorIs(t, err, shared.ErrNotFound)
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestStorage_ListLinksByUserID(t *testing.T) {
	store := New()
	ctx := context.Background()

	records := []models.URLRecord{
		{ID: "short1", URL: "https://example.com", UserID: "user1"},
		{ID: "short2", URL: "https://test.com", UserID: "user1"},
		{ID: "short3", URL: "https://another.com", UserID: "user2"},
	}

	for _, r := range records {
		require.NoError(t, store.Put(ctx, r))
	}

	testCases := []struct {
		wantErr error
		name    string
		userID  string
		baseURL string
		wantLen int
	}{
		{nil, "List for user1", "user1", "https://short.ly", 2},
		{nil, "List for user2", "user2", "https://short.ly", 1},
		{shared.ErrNotFound, "List for non-existent user", "user3", "https://short.ly", 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			links, err := store.ListLinksByUserID(ctx, tc.baseURL, tc.userID)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
				assert.Len(t, links, tc.wantLen)
			}
		})
	}
}

func TestStorage_PutBatch(t *testing.T) {
	store := New()
	ctx := context.Background()

	testCases := []struct {
		wantErr error
		name    string
		records []models.URLRecord
	}{
		{
			name: "Insert batch of records",
			records: []models.URLRecord{
				{ID: "batch1", URL: "https://batch1.com", UserID: "user1"},
				{ID: "batch2", URL: "https://batch2.com", UserID: "user1"},
			},
			wantErr: nil,
		},
		{
			name:    "Empty batch",
			records: []models.URLRecord{},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := store.PutBatch(ctx, tc.records)
			require.NoError(t, err)

			for _, r := range tc.records {
				gotURL, err := store.Get(ctx, r.ID)
				require.NoError(t, err)
				assert.Equal(t, r.URL, gotURL)
			}
		})
	}
}

func TestStorage_PingAndClose(t *testing.T) {
	store := New()

	testCases := []struct {
		fn   func() error
		name string
	}{
		{name: "Close", fn: store.Close},
		{name: "Ping", fn: store.Ping},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NoError(t, tc.fn())
		})
	}
}

func Test_Stats(t *testing.T) {
	im := New()

	records := []models.URLRecord{
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

	for i, rec := range records {
		t.Run(fmt.Sprintf("put #%d", i), func(t *testing.T) {
			err := im.Put(context.Background(), rec)
			require.NoError(t, err)
		})
	}

	t.Run("check stats", func(t *testing.T) {
		stats, err := im.Stats(context.Background())
		require.NoError(t, err)

		assert.Equal(t, 15, stats.Urls, "URL count mismatch")
		assert.Equal(t, 2, stats.Users, "User count mismatch")
	})
}

package storages_test

import (
	"context"
	"fmt"
	"log"

	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/inmem"
)

func ExampleStorage_Put() {
	ctx := context.Background()
	storage := inmem.New()

	record := models.URLRecord{
		ID:     "short123",
		URL:    "https://example.com",
		UserID: "user123",
	}

	err := storage.Put(ctx, record)
	if err != nil {
		log.Fatalf("failed to put URL: %v", err)
	}

	fmt.Println("URL inserted successfully")
	// Output: URL inserted successfully
}

func ExampleStorage_Get() {
	ctx := context.Background()
	storage := inmem.New()

	record := models.URLRecord{
		ID:     "short123",
		URL:    "https://example.com",
		UserID: "user123",
	}
	// Store the record in the in-memory storage
	err := storage.Put(ctx, record)
	if err != nil {
		log.Fatalf("failed to put URL: %v", err)
	}

	url, err := storage.Get(ctx, "short123")
	if err != nil {
		log.Fatalf("failed to get URL: %v", err)
	}

	fmt.Println("Retrieved URL:", url)
	// Output: Retrieved URL: https://example.com
}

func ExampleStorage_ListLinksByUserID() {
	ctx := context.Background()
	storage := inmem.New()

	record := models.URLRecord{
		ID:     "short123",
		URL:    "http://short.url/short123",
		UserID: "user123",
	}
	// Store the record in the in-memory storage
	err := storage.Put(ctx, record)
	if err != nil {
		log.Fatalf("failed to put URL: %v", err)
	}

	records, err := storage.ListLinksByUserID(ctx, "http://short.url", "user123")
	if err != nil {
		log.Fatalf("failed to list URLs: %v", err)
	}

	for _, record := range records {
		fmt.Println("URL:", record.URL)
	}
	// Output: URL: http://short.url/short123
}

func ExampleStorage_DeleteUserURLs() {
	ctx := context.Background()
	storage := inmem.New()

	err := storage.DeleteUserURLs(ctx, []string{"short123"}, "user123")
	if err != nil {
		log.Fatalf("failed to delete URLs: %v", err)
	}

	fmt.Println("URLs deleted successfully")
	// Output: URLs deleted successfully
}

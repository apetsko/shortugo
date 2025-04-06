package infile

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"slices"
	"sync"

	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/shared"
)

// File permissions for user read/write, group read, others read.
const FilePermUserRWGroupROthersR = 0644

// CustomBool is a custom boolean type for JSON marshaling/unmarshaling.
type CustomBool bool

// UnmarshalJSON unmarshals a boolean from JSON.
func (b *CustomBool) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		*b = false
		return nil
	}

	if data[0] == '1' {
		*b = true
	} else {
		*b = false
	}
	return nil
}

// MarshalJSON marshals a boolean to JSON.
func (b CustomBool) MarshalJSON() ([]byte, error) {
	if b {
		return []byte("1"), nil
	}
	return []byte("0"), nil
}

// Storage represents a storage backed by a file.
type Storage struct {
	file    *os.File
	encoder *json.Encoder
	mu      sync.Mutex
}

// New creates a new Storage instance with the given filename.
func New(filename string) (*Storage, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, FilePermUserRWGroupROthersR)
	if err != nil {
		return nil, err
	}

	return &Storage{
		file:    f,
		encoder: json.NewEncoder(f),
	}, nil
}

// Close closes the storage file.
func (f *Storage) Close() error {
	return f.file.Close()
}

// Put stores a URLRecord in the storage.
func (f *Storage) Put(ctx context.Context, r models.URLRecord) (err error) {
	if err := ctx.Err(); err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if err := f.encoder.Encode(r); err != nil {
		return err
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	if err := f.file.Sync(); err != nil {
		return fmt.Errorf("error sync file: %w", err)
	}

	return nil
}

// PutBatch stores multiple URLRecords in the storage.
func (f *Storage) PutBatch(ctx context.Context, rr []models.URLRecord) (err error) {
	for _, r := range rr {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := f.encoder.Encode(r); err != nil {
			return err
		}

		if err := f.file.Sync(); err != nil {
			return fmt.Errorf("error sync file: %w", err)
		}
	}
	return nil
}

// Get retrieves the original URL for a given short URL.
func (f *Storage) Get(ctx context.Context, shortURL string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	if _, err := f.file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("error setting file seek: %w", err)
	}

	r := new(models.URLRecord)
	scanner := bufio.NewScanner(f.file)
	for scanner.Scan() {
		data := scanner.Bytes()
		err := json.Unmarshal(data, &r)
		if err != nil {
			return "", fmt.Errorf("failed unmarshal: %w", err)
		}
		fmt.Println("r", r)
		if r.ID == shortURL {
			if r.Deleted {
				return "", errors.New(http.StatusText(http.StatusGone))
			}

			return r.URL, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	return "", fmt.Errorf("URL not found: %s. %w", shortURL, shared.ErrNotFound)
}

// ListLinksByUserID lists all URLs associated with a user ID.
func (f *Storage) ListLinksByUserID(ctx context.Context, baseURL, userID string) ([]models.URLRecord, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if _, err := f.file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("error setting file seek: %w", err)
	}

	rr := make([]models.URLRecord, 0)
	r := new(models.URLRecord)
	scanner := bufio.NewScanner(f.file)
	for scanner.Scan() {
		data := scanner.Bytes()
		if err := json.Unmarshal(data, r); err != nil {
			return nil, err
		}

		if r.UserID == userID && !r.Deleted {
			r.ID = baseURL + "/" + r.ID
			rr = append(rr, *r)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if len(rr) == 0 {
		return rr, fmt.Errorf("URLs not found for UserID: %s. %w", userID, shared.ErrNotFound)
	}
	return rr, nil
}

// DeleteUserURLs deletes multiple URLs associated with a user ID.
func (f *Storage) DeleteUserURLs(ctx context.Context, ids []string, userID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	tmpFilename := f.file.Name() + ".tmp"
	tmpFile, err := os.OpenFile(tmpFilename, os.O_CREATE|os.O_WRONLY, FilePermUserRWGroupROthersR)
	if err != nil {
		return fmt.Errorf("error creating temp file: %w", err)
	}
	defer os.Remove(tmpFilename)
	defer tmpFile.Close()

	if err := f.copyAndMarkDeleted(tmpFile, ids, userID); err != nil {
		return err
	}

	if err := f.replaceFile(tmpFilename); err != nil {
		return err
	}

	return nil
}

// copyAndMarkDeleted copies records to a temporary file and marks specified records as deleted.
func (f *Storage) copyAndMarkDeleted(tmpFile *os.File, ids []string, userID string) error {
	if _, err := f.file.Seek(0, 0); err != nil {
		return fmt.Errorf("error setting file seek: %w", err)
	}

	scanner := bufio.NewScanner(f.file)
	writer := bufio.NewWriter(tmpFile)

	for scanner.Scan() {
		r, err := f.parseRecord(scanner.Bytes())
		if err != nil {
			return err
		}

		if shouldDelete(r, ids, userID) {
			r.Deleted = true
		}

		if err := writeRecord(writer, r); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return writer.Flush()
}

// shouldDelete determines if a record should be marked as deleted.
func shouldDelete(r *models.URLRecord, ids []string, userID string) bool {
	return r.UserID == userID && slices.Contains(ids, r.ID) && !r.Deleted
}

// parseRecord parses a URLRecord from a byte slice.
func (f *Storage) parseRecord(line []byte) (*models.URLRecord, error) {
	r := new(models.URLRecord)
	if err := json.Unmarshal(line, r); err != nil {
		return nil, err
	}
	return r, nil
}

// writeRecord writes a URLRecord to a buffered writer.
func writeRecord(writer *bufio.Writer, r *models.URLRecord) error {
	newLine, err := json.Marshal(r)
	if err != nil {
		return err
	}

	if _, err := writer.Write(newLine); err != nil {
		return fmt.Errorf("error writing updated record to temp file: %w", err)
	}

	if _, err := writer.WriteString("\n"); err != nil {
		return fmt.Errorf("error writing newline to temp file: %w", err)
	}

	return nil
}

// replaceFile replaces the original storage file with the temporary file.
func (f *Storage) replaceFile(tmpFilename string) error {
	f.file.Close()

	if err := os.Rename(tmpFilename, f.file.Name()); err != nil {
		return fmt.Errorf("error replacing storage file: %w", err)
	}

	var err error
	f.file, err = os.OpenFile(f.file.Name(), os.O_RDWR|os.O_CREATE|os.O_APPEND, FilePermUserRWGroupROthersR)
	if err != nil {
		return fmt.Errorf("error reopening storage file: %w", err)
	}
	f.encoder = json.NewEncoder(f.file)

	return nil
}

// Ping checks the storage health.
func (f *Storage) Ping() error {
	return nil
}

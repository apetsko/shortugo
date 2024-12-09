package utils

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateURLS(count int) []string {

	urls := make([]string, count)

	protocols := []string{
		"http://",
		"https://",
		"ftp://",
		"ftps://",
		"sftp://",
		"mailto://",
		"telnet://",
		"file://",
		"data://",
		"ws://",
		"wss://",
		"bluetooth://",
	}

	for i := range urls {
		u := fmt.Sprintf("%s%s.%s/%s/%s",
			protocols[seededRand.Intn(len(protocols))], //protocol
			randomString(seededRand.Intn(10)+1),        //2domain
			randomString(seededRand.Intn(3)+2),         //1domain
			randomString(seededRand.Intn(9)+2),         //path1
			randomString(seededRand.Intn(13)+4),        //path2
		)
		urls[i] = u
	}

	return urls
}

func randomString(lenght int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, lenght)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func Test_Generate(t *testing.T) {
	m := make(map[string]string)
	urls := generateURLS(1000_000)
	for i, u := range urls {
		t.Run(fmt.Sprintf("URL #%d", i), func(t *testing.T) {
			ID := Generate(u)
			if ur, ok := m[ID]; ok && ur != u {
				t.Errorf("already has same ID %q with URL: %q, new URL: %q", ID, m[ID], u)
			}
			m[ID] = u
		})
	}
}

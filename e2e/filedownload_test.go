package e2e

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
)

func NewFileRequest(filename string) (*http.Request, error) {
	// Use Config to build URL
	return http.NewRequest("GET", Config.GetServerURL("/files/"+filename), nil)
}

func NewFileRequestWithMethod(method string, filename string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, Config.GetServerURL("/files/"+filename), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return req, nil
}

func ExecuteRequest(req *http.Request) (*http.Response, error) {
	// Add test name as request ID if we're in a test context
	if testName := req.Context().Value("testName"); testName != nil {
		req.Header.Set("X-Test-ID", testName.(string))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

func TestFiles(t *testing.T) {
	t.Run("Download existing file returns 200 and correct content", func(t *testing.T) {
		req, err := NewFileRequest("foo")
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		// Add test name to request context
		ctx := context.WithValue(req.Context(), "testName", t.Name())
		req = req.WithContext(ctx)

		resp, err := ExecuteRequest(req)
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}
		defer resp.Body.Close()

		// Check response status
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got: %d", resp.StatusCode)
		}

		// Check Content-Type header
		contentType := resp.Header.Get("Content-Type")
		if contentType != "application/octet-stream" {
			t.Errorf("Expected Content-Type 'application/octet-stream', got: '%s'", contentType)
		}

		// Read response body into memory
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		// Verify content
		expectedContent := "Hello, World!"
		if strings.TrimSpace(string(content)) != expectedContent {
			t.Errorf("Expected content '%s', got: '%s'", expectedContent, string(content))
		}
	})

	t.Run("Download non-existent file returns 404", func(t *testing.T) {
		req, err := NewFileRequest("non_existant_file")
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := ExecuteRequest(req)
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}
		defer resp.Body.Close()

		// Check response status
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got: %d", resp.StatusCode)
		}
	})

	t.Run("Upload new file returns 201 and file is readable", func(t *testing.T) {
		t.Cleanup(func() {
			cleanupTestFiles(t)
		})

		content := "12345"
		filename := "test-upload.txt"

		req, err := NewFileRequestWithMethod("POST", filename, strings.NewReader(content))
		if err != nil {
			t.Fatalf("Failed to create upload request: %v", err)
		}

		req.Header.Set("Content-Type", "application/octet-stream")

		resp, err := ExecuteRequest(req)
		if err != nil {
			t.Fatalf("Failed to upload file: %v", err)
		}
		defer resp.Body.Close()

		// Check upload response
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Expected status 201, got: %d", resp.StatusCode)
		}

		req, err = NewFileRequest("test-upload.txt")
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err = ExecuteRequest(req)
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}
		defer resp.Body.Close()

		// Check download response
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got: %d", resp.StatusCode)
		}

		// Check Content-Type header
		contentType := resp.Header.Get("Content-Type")
		if contentType != "application/octet-stream" {
			t.Errorf("Expected Content-Type 'application/octet-stream', got: '%s'", contentType)
		}

		// Read and verify content
		downloadedContent, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		if string(downloadedContent) != content {
			t.Errorf("Expected content '%s', got: '%s'", content, string(downloadedContent))
		}
	})
}

func cleanupTestFiles(t *testing.T) {
	// List of test files to cleanup
	testFiles := []string{
		"test-upload.txt",
	}

	for _, filename := range testFiles {
		fullPath := path.Join(Config.Directory, filename)
		err := os.Remove(fullPath)
		if err != nil && !os.IsNotExist(err) {
			t.Logf("Failed to cleanup file %s: %v", filename, err)
		}
	}
}

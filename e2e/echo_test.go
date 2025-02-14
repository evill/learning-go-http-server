package e2e

import (
	"io"
	"net/http"
	"testing"
)

func TestEcho(t *testing.T) {
	t.Run("Echo endpoint returns content after /echo/", func(t *testing.T) {
		req, err := http.NewRequest("GET", Config.GetServerURL("/echo/hello"), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := ExecuteRequest(req)
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		expected := "hello"
		if string(body) != expected {
			t.Errorf("Expected body '%s', got: '%s'", expected, string(body))
		}
	})

	t.Run("Root path returns 200 OK", func(t *testing.T) {
		req, err := http.NewRequest("GET", Config.GetServerURL("/"), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := ExecuteRequest(req)
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got: %d", resp.StatusCode)
		}
	})
}

func TestGetEchoWithGzip(t *testing.T) {
	req, err := http.NewRequest("GET", Config.GetServerURL("/echo/foo"), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Test-ID", t.Name())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got: %d", resp.StatusCode)
	}

	// Read decompressed content
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if string(body) != "foo" {
		t.Errorf("Expected body 'foo', got: %s", string(body))
	}
}

func TestGetEchoWithoutGzip(t *testing.T) {
	echoContent := "foo"
	req, err := http.NewRequest("GET", Config.GetServerURL("/echo/"+echoContent), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Test-ID", t.Name())
	// Explicitly request uncompressed content
	req.Header.Set("Accept-Encoding", "identity")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got: %d", resp.StatusCode)
	}

	// Read decompressed content
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if string(body) != echoContent {
		t.Errorf("Expected body 'foo', got: %s", string(body))
	}
}

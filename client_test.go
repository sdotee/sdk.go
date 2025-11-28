package seesdk

import (
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	baseUrl := DEFAULT_BASE_URL

	if os.Getenv("SEE_API_KEY") == "" {
		t.Skip("SEE_API_KEY not set, skipping integration test")
	}

	if os.Getenv("SEE_BASE_URL") != "" {
		baseUrl = os.Getenv("SEE_BASE_URL")
	}

	client := NewClient(Config{
		BaseURL: baseUrl,
		APIKey:  os.Getenv("SEE_API_KEY"),
	})

	if client == nil {
		t.Fatal("Expected client to be created")
	}

	domains, err := client.GetDomains()
	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}

	if len(domains.Data.Domains) == 0 {
		t.Fatal("Expected at least one domain, got zero")
	}

	tags, err := client.GetTags()
	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}

	if tags.Code != 200 {
		t.Fatal("Expected response code 200, got:", tags.Code)
	}

	response, err := client.CreateShortURL(CreateShortURLRequest{
		Domain:    "a.see-test.com",
		TargetURL: "https://www.google.com/",
	})

	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}

	if response.Code != 200 {
		t.Errorf("Expected response code 200, got: %d", response.Code)
	}

	result, err := client.UpdateShortURL(UpdateShortURLRequest{
		Domain:    "a.see-test.com",
		Slug:      response.Data.Slug,
		Title:     "Google",
		TargetUrl: "https://www.google.com/search?q=see+sdk",
	})

	if err != nil {
		t.Fatal("Expected no error on delete, got:", err)
	}

	if result.Code != 200 {
		t.Errorf("Expected delete response code 200, got: %d", result.Code)
	}

	result2, err := client.DeleteShortURL(DeleteURLRequest{
		Domain: "a.see-test.com",
		Slug:   response.Data.Slug,
	})

	if err != nil {
		t.Fatal("Expected no error on delete, got:", err)
	}

	if result2.Code != 200 {
		t.Errorf("Expected delete response code 200, got: %d", result.Code)
	}
}

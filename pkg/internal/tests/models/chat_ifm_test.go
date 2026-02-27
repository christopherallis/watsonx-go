package test

import (
	"context"
	"os"
	"testing"

	wx "github.com/IBM/watsonx-go/pkg/models"
)

// TestChatIFMLite validates chat functionality against an IFM Lite (CPD) deployment.
// Required env vars:
//
//	CPD_USERNAME  - CPD username
//	API_KEY       - CPD API key
//	MODEL_ID      - model ID to test (e.g. "ibm/granite-3-8b-instruct")
//	WATSONX_URL   - full WatsonX URL (e.g. "https://cpd.example.com")
//
// No project_id or space_id is used — IFM Lite does not require one.
func TestChatIFMLite(t *testing.T) {
	cpdUsername := os.Getenv("CPD_USERNAME")
	apiKey := os.Getenv("API_KEY")
	modelID := os.Getenv("MODEL_ID")
	watsonxURL := os.Getenv("WATSONX_URL")

	if cpdUsername == "" {
		t.Fatal("CPD_USERNAME env var is required")
	}
	if apiKey == "" {
		t.Fatal("API_KEY env var is required")
	}
	if modelID == "" {
		t.Fatal("MODEL_ID env var is required")
	}
	if watsonxURL == "" {
		t.Fatal("WATSONX_URL env var is required")
	}

	client, err := wx.NewClient(
		wx.WithCPD(watsonxURL, cpdUsername),
		wx.WithCPDAPIKey(apiKey),
	)
	if err != nil {
		t.Fatalf("Failed to create IFM Lite client: %v", err)
	}

	ctx := context.Background()

	messages := []wx.ChatMessage{
		wx.CreateUserMessage("What is the capital of France?"),
	}

	response, err := client.Chat(ctx, modelID, messages)
	if err != nil {
		t.Fatalf("Chat request failed: %v", err)
	}

	if len(response.Choices) == 0 {
		t.Fatal("Expected at least one choice in response")
	}

	if response.Choices[0].Message == nil {
		t.Fatal("Expected message in first choice")
	}

	content := response.Choices[0].Message.Content.GetText()
	if content == "" {
		t.Fatal("Expected non-empty response content")
	}

	t.Logf("Model: %s", response.ModelID)
	t.Logf("Response: %s", content)
}

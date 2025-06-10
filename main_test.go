package main

import (
	"context"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestCaffeinateServer(t *testing.T) {
	srv := NewCaffeinateServer()

	t.Run("Start caffeinate with idle flag", func(t *testing.T) {
		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "caffeinate_start",
				Arguments: map[string]interface{}{
					"idle": true,
				},
			},
		}

		result, err := srv.handleCaffeinateStart(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result == nil {
			t.Fatal("expected result, got nil")
		}

		// Check if result contains expected text
		if len(result.Content) == 0 {
			t.Fatal("expected content in result")
		}

		// Check if it's not an error
		if result.IsError {
			t.Fatal("expected success result, got error")
		}

		// Extract text content and verify it contains expected information
		content := getTextContent(result)
		if !strings.Contains(content, "Started caffeinate") {
			t.Errorf("expected result to contain 'Started caffeinate', got: %s", content)
		}
		if !strings.Contains(content, "-i") {
			t.Errorf("expected result to contain '-i' flag, got: %s", content)
		}
	})

	t.Run("List caffeinate processes", func(t *testing.T) {
		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name:      "caffeinate_list",
				Arguments: map[string]interface{}{},
			},
		}

		result, err := srv.handleCaffeinateList(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result == nil {
			t.Fatal("expected result, got nil")
		}

		// Check if result contains content
		if len(result.Content) == 0 {
			t.Fatal("expected content in result")
		}

		// Should not be an error
		if result.IsError {
			t.Fatal("expected success result, got error")
		}
	})

	t.Run("Stop non-existent caffeinate process", func(t *testing.T) {
		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name: "caffeinate_stop",
				Arguments: map[string]interface{}{
					"id": "non-existent-id",
				},
			},
		}

		result, err := srv.handleCaffeinateStop(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result == nil {
			t.Fatal("expected result, got nil")
		}

		// Should return an error result for non-existent process
		if !result.IsError {
			t.Fatal("expected error result for non-existent process")
		}

		// Check error message
		content := getTextContent(result)
		if !strings.Contains(content, "No caffeinate process found") {
			t.Errorf("expected error message about non-existent process, got: %s", content)
		}
	})

	// Cleanup any active processes
	srv.mu.Lock()
	for _, cmd := range srv.activeProcesses {
		cmd.Process.Kill()
	}
	srv.mu.Unlock()
}

// Helper function to extract text content from CallToolResult
func getTextContent(result *mcp.CallToolResult) string {
	var texts []string
	for _, content := range result.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			texts = append(texts, textContent.Text)
		}
	}
	return strings.Join(texts, " ")
}

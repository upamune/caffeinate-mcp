package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// Example of creating a CallToolRequest for testing
func ExampleCallToolRequest() {
	// Create a CallToolRequest with map arguments
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "caffeinate_start",
			Arguments: map[string]interface{}{
				"idle":    true,
				"display": true,
				"timeout": 300,
			},
		},
	}

	// Access arguments using helper methods
	idle := req.GetBool("idle", false)
	timeout := req.GetInt("timeout", 0)

	fmt.Printf("Idle: %v, Timeout: %d\n", idle, timeout)
	// Output: Idle: true, Timeout: 300
}

// Example of checking CallToolResult
func ExampleCallToolResult() {
	// Create a sample result (usually returned from a tool handler)
	result := mcp.NewToolResultText("Started caffeinate with ID: 12345")

	// Check if it's an error
	if result.IsError {
		fmt.Println("Error occurred")
		return
	}

	// Extract text content
	for _, content := range result.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			fmt.Printf("Result: %s\n", textContent.Text)
		}
	}
	// Output: Result: Started caffeinate with ID: 12345
}

// Example of a tool handler for testing
func Example_toolHandler() {
	handler := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get required string parameter
		name, err := req.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError("Name is required"), nil
		}

		// Get optional parameters with defaults
		age := req.GetInt("age", 0)
		isActive := req.GetBool("active", true)

		// Return success result
		return mcp.NewToolResultText(fmt.Sprintf("Hello %s, age: %d, active: %v", name, age, isActive)), nil
	}

	// Test the handler
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "greet",
			Arguments: map[string]interface{}{
				"name": "Alice",
				"age":  25,
			},
		},
	}

	result, _ := handler(context.Background(), req)
	for _, content := range result.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			fmt.Println(textContent.Text)
		}
	}
	// Output: Hello Alice, age: 25, active: true
}

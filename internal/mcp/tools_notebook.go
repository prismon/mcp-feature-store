package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/prismon/synthesis/internal/domain"
)

// createLibraryTool defines the create_library tool
func createLibraryTool() mcp.Tool {
	return mcp.Tool{
		Name:        "create_library",
		Description: "Create a new library in a tenant",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"tenantId": map[string]interface{}{
					"type":        "string",
					"description": "ID of the tenant",
				},
				"libraryId": map[string]interface{}{
					"type":        "string",
					"description": "Unique identifier for the library",
				},
				"owner": map[string]interface{}{
					"type":        "string",
					"description": "Email address of the library owner",
				},
				"display_name": map[string]interface{}{
					"type":        "string",
					"description": "Display name for the library",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Description of the library",
				},
			},
			Required: []string{"tenantId", "libraryId", "owner", "display_name"},
		},
	}
}

// handleCreateLibrary handles the create_library tool invocation
func (s *Server) handleCreateLibrary(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Placeholder implementation
	return mcp.NewToolResultText(`{"message": "Library creation not yet implemented"}`), nil
}

// createNotebookTool defines the create_notebook tool
func createNotebookTool() mcp.Tool {
	return mcp.Tool{
		Name:        "create_notebook",
		Description: "Create a new notebook in a library",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"tenantId": map[string]interface{}{
					"type":        "string",
					"description": "ID of the tenant",
				},
				"libraryId": map[string]interface{}{
					"type":        "string",
					"description": "ID of the library",
				},
				"notebookId": map[string]interface{}{
					"type":        "string",
					"description": "Unique identifier for the notebook",
				},
				"display_name": map[string]interface{}{
					"type":        "string",
					"description": "Display name for the notebook",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Description of the notebook",
				},
				"initial_markdown": map[string]interface{}{
					"type":        "string",
					"description": "Initial markdown content",
				},
			},
			Required: []string{"tenantId", "libraryId", "notebookId", "display_name"},
		},
	}
}

// handleCreateNotebook handles the create_notebook tool invocation
func (s *Server) handleCreateNotebook(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse arguments
	var args struct {
		TenantID        string `json:"tenantId"`
		LibraryID       string `json:"libraryId"`
		NotebookID      string `json:"notebookId"`
		DisplayName     string `json:"display_name"`
		Description     string `json:"description"`
		InitialMarkdown string `json:"initial_markdown"`
	}

	argsBytes, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal arguments: %v", err)), nil
	}

	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse arguments: %v", err)), nil
	}

	// Create notebook
	notebook := &domain.Notebook{
		TenantID:    args.TenantID,
		ID:          args.NotebookID,
		LibraryID:   args.LibraryID,
		Status:      "draft",
		Owner:       "system", // TODO: Get from auth context
		DisplayName: args.DisplayName,
		Description: args.Description,
		Contents: domain.NotebookContents{
			Data: domain.NotebookData{
				Markdown: args.InitialMarkdown,
			},
			ContentBlocks: []domain.ContentBlock{},
		},
		Notifications: []domain.Notification{},
	}

	if err := s.notebookRepo.Create(ctx, notebook); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create notebook: %v", err)), nil
	}

	// Return success response
	result := map[string]interface{}{
		"notebookUri": notebook.URI(),
		"notebookId":  notebook.ID,
		"message":     "Notebook created successfully",
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(resultBytes)), nil
}

// appendContentBlockTool defines the append_content_block tool
func appendContentBlockTool() mcp.Tool {
	return mcp.Tool{
		Name:        "append_content_block",
		Description: "Append a content block to a notebook",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"notebookId": map[string]interface{}{
					"type":        "string",
					"description": "ID of the notebook",
				},
				"content_type": map[string]interface{}{
					"type":        "string",
					"description": "Content type (e.g., text/markdown, binary/image)",
				},
				"data": map[string]interface{}{
					"type":        "string",
					"description": "Content data",
				},
				"parent_uid": map[string]interface{}{
					"type":        "string",
					"description": "UID of parent block (optional)",
				},
			},
			Required: []string{"notebookId", "content_type", "data"},
		},
	}
}

// handleAppendContentBlock handles the append_content_block tool invocation
func (s *Server) handleAppendContentBlock(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Placeholder implementation
	return mcp.NewToolResultText(`{"message": "Append content block not yet implemented"}`), nil
}

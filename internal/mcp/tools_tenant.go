package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/prismon/synthesis/internal/domain"
)

// createTenantTool defines the create_tenant tool
func createTenantTool() mcp.Tool {
	return mcp.Tool{
		Name:        "create_tenant",
		Description: "Create a new tenant in Synthesis",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"tenantId": map[string]interface{}{
					"type":        "string",
					"description": "Unique identifier for the tenant",
				},
				"owner": map[string]interface{}{
					"type":        "string",
					"description": "Email address of the tenant owner",
				},
				"display_name": map[string]interface{}{
					"type":        "string",
					"description": "Display name for the tenant",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Description of the tenant",
				},
				"labels": map[string]interface{}{
					"type":        "object",
					"description": "Key-value labels for the tenant",
				},
			},
			Required: []string{"tenantId", "owner", "display_name"},
		},
	}
}

// handleCreateTenant handles the create_tenant tool invocation
func (s *Server) handleCreateTenant(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse arguments
	var args struct {
		TenantID    string            `json:"tenantId"`
		Owner       string            `json:"owner"`
		DisplayName string            `json:"display_name"`
		Description string            `json:"description"`
		Labels      map[string]string `json:"labels"`
	}

	argsBytes, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal arguments: %v", err)), nil
	}

	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse arguments: %v", err)), nil
	}

	// Validate required fields
	if args.TenantID == "" || args.Owner == "" || args.DisplayName == "" {
		return mcp.NewToolResultError("tenantId, owner, and display_name are required"), nil
	}

	// Create tenant
	tenant := &domain.Tenant{
		ID:           args.TenantID,
		Owner:        args.Owner,
		DisplayName:  args.DisplayName,
		Description:  args.Description,
		Labels:       args.Labels,
		Version:      "1",
		LastModified: time.Now(),
	}

	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create tenant: %v", err)), nil
	}

	// Return success response
	result := map[string]interface{}{
		"tenantUri": tenant.URI(),
		"tenantId":  tenant.ID,
		"message":   "Tenant created successfully",
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(resultBytes)), nil
}

// getTenantTool defines the get_tenant tool
func getTenantTool() mcp.Tool {
	return mcp.Tool{
		Name:        "get_tenant",
		Description: "Get tenant information by ID",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"tenantId": map[string]interface{}{
					"type":        "string",
					"description": "Unique identifier for the tenant",
				},
			},
			Required: []string{"tenantId"},
		},
	}
}

// handleGetTenant handles the get_tenant tool invocation
func (s *Server) handleGetTenant(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse arguments
	var args struct {
		TenantID string `json:"tenantId"`
	}

	argsBytes, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal arguments: %v", err)), nil
	}

	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse arguments: %v", err)), nil
	}

	// Get tenant
	tenant, err := s.tenantRepo.Get(ctx, args.TenantID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get tenant: %v", err)), nil
	}

	// Marshal tenant to JSON
	resultBytes, err := json.Marshal(tenant)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal tenant: %v", err)), nil
	}

	return mcp.NewToolResultText(string(resultBytes)), nil
}

package mcp

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/prismon/synthesis/internal/postgres"
)

// Server represents the Synthesis MCP server
type Server struct {
	mcpServer    *server.MCPServer
	tenantRepo   *postgres.TenantRepository
	notebookRepo *postgres.NotebookRepository
	vectorRepo   *postgres.VectorRepository
	graphRepo    *postgres.GraphRepository
}

// NewServer creates a new MCP server
func NewServer(db *postgres.DB) (*Server, error) {
	s := &Server{
		tenantRepo:   postgres.NewTenantRepository(db),
		notebookRepo: postgres.NewNotebookRepository(db),
		vectorRepo:   postgres.NewVectorRepository(db),
		graphRepo:    postgres.NewGraphRepository(db),
	}

	// Create MCP server with capabilities
	mcpServer := server.NewMCPServer(
		"Synthesis MCP Server",
		"1.0.0",
		server.WithResourceCapabilities(true, true), // supports list and read
		server.WithToolCapabilities(true),           // supports tools
	)

	s.mcpServer = mcpServer

	// Register tools
	s.registerTools()

	return s, nil
}

// registerTools registers all MCP tools
func (s *Server) registerTools() {
	// Tenant tools
	s.mcpServer.AddTool(createTenantTool(), s.handleCreateTenant)
	s.mcpServer.AddTool(getTenantTool(), s.handleGetTenant)

	// Library tools
	s.mcpServer.AddTool(createLibraryTool(), s.handleCreateLibrary)

	// Notebook tools
	s.mcpServer.AddTool(createNotebookTool(), s.handleCreateNotebook)
	s.mcpServer.AddTool(appendContentBlockTool(), s.handleAppendContentBlock)

	// Search tools
	s.mcpServer.AddTool(semanticSearchNotebooksTool(), s.handleSemanticSearchNotebooks)
	s.mcpServer.AddTool(graphQueryResourcesTool(), s.handleGraphQueryResources)
}

// Start starts the MCP server with stdio transport
func (s *Server) Start() error {
	// Create stdio server
	stdioServer := server.NewStdioServer(s.mcpServer)

	// Start listening on stdio
	ctx := context.Background()
	return stdioServer.Listen(ctx, os.Stdin, os.Stdout)
}

// Example resource handler (not used but shows pattern)
func (s *Server) handleListResources(ctx context.Context) ([]mcp.Resource, error) {
	resources := []mcp.Resource{}

	// List all tenants and create resource entries
	tenants, err := s.tenantRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}

	for _, tenant := range tenants {
		resources = append(resources, mcp.Resource{
			URI:         tenant.URI(),
			Name:        tenant.ID,
			Description: tenant.Description,
			MIMEType:    "application/json",
		})
	}

	return resources, nil
}

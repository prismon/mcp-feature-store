# Synthesis - MCP Feature Store

Synthesis is a tool for AI systems that acts as a **central nervous system** for individual users and small to mid-sized companies. This implementation is built as an **MCP server in Go** backed by **PostgreSQL** with vector and graph capabilities.

## Features

- **MCP Server**: Built with `github.com/mark3labs/mcp-go` for AI agent integration
- **PostgreSQL Backend**: Enhanced with pgvector for embeddings and Apache AGE for graph queries
- **Domain Model**: Supports Tenants, Libraries, Notebooks, Features, Types, Products, and Tools
- **Vector Search**: Semantic search capabilities using pgvector
- **Graph Queries**: Relationship traversal using Apache AGE
- **REST API**: Traditional HTTP API for non-MCP clients
- **Docker Support**: Full containerization with docker-compose

## Architecture

```
synthesis/
├── cmd/
│   ├── synthesis-mcp/     # MCP server (stdio transport)
│   └── synthesis-api/     # REST API server
├── internal/
│   ├── config/            # Configuration management
│   ├── domain/            # Domain models (Tenant, Library, Notebook, etc.)
│   ├── postgres/          # Database layer (repositories, migrations)
│   ├── mcp/               # MCP server implementation
│   └── rest/              # REST API handlers
├── db/
│   └── migrations/        # SQL migration scripts
├── deployments/
│   ├── Dockerfile.postgres  # PostgreSQL with pgvector + AGE
│   └── init-extensions.sql  # Extension initialization
├── Dockerfile             # Application Dockerfile
└── docker-compose.yml     # Local development stack
```

## Technology Stack

- **Language**: Go 1.22
- **MCP Library**: github.com/mark3labs/mcp-go
- **Database**: PostgreSQL 16
  - pgvector: Vector embeddings and similarity search
  - Apache AGE: Graph database capabilities
- **HTTP Framework**: Gorilla Mux
- **Auth (Planned)**: OIDC + OPA

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.22+ (for local development)

### Using Docker Compose (Recommended)

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd mcp-feature-store
   ```

2. **Start the stack**:
   ```bash
   docker-compose up -d
   ```

   This will start:
   - PostgreSQL with pgvector and Apache AGE (port 5432)
   - Synthesis MCP Server (stdio transport)
   - Synthesis REST API (port 8080)
   - Open Policy Agent (port 8181)

3. **Check the logs**:
   ```bash
   docker-compose logs -f synthesis-api
   docker-compose logs -f synthesis-mcp
   ```

4. **Test the REST API**:
   ```bash
   curl http://localhost:8080/health
   ```

### Building Locally

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Build binaries**:
   ```bash
   go build -o bin/synthesis-mcp ./cmd/synthesis-mcp
   go build -o bin/synthesis-api ./cmd/synthesis-api
   ```

3. **Start PostgreSQL** (using docker):
   ```bash
   docker-compose up -d postgres
   ```

4. **Run the servers**:
   ```bash
   # Terminal 1 - MCP Server
   export DATABASE_URL="postgres://synthesis:synthesis_dev_password@localhost:5432/synthesis?sslmode=disable"
   ./bin/synthesis-mcp

   # Terminal 2 - REST API
   export DATABASE_URL="postgres://synthesis:synthesis_dev_password@localhost:5432/synthesis?sslmode=disable"
   ./bin/synthesis-api
   ```

## Configuration

Configuration is done via environment variables:

### Database
- `DATABASE_URL`: PostgreSQL connection string
- `DB_MAX_OPEN_CONNS`: Maximum open connections (default: 25)
- `DB_MAX_IDLE_CONNS`: Maximum idle connections (default: 5)
- `DB_CONN_MAX_LIFETIME`: Connection max lifetime in seconds (default: 300)

### MCP Server
- `MCP_HTTP_PORT`: HTTP port for MCP server (default: 8081)
- `LOG_LEVEL`: Log level (default: info)

### REST API
- `API_PORT`: REST API port (default: 8080)
- `WS_PORT`: WebSocket port (default: 8082)

### Security (Planned)
- `OIDC_ISSUER_URL`: OIDC issuer URL
- `OIDC_CLIENT_ID`: OIDC client ID
- `OPA_URL`: OPA server URL (default: http://localhost:8181)

## Domain Model

### Tenant
Top-level organization unit with owner, display name, labels, and version control.

### Library
Collection of notebooks within a tenant.

### Notebook
Primary editable element with:
- Markdown content
- Hierarchical content blocks (text, images, etc.)
- Status tracking (draft, approved)
- Notifications

### Feature
Derived data associated with resources:
- Text, JSON, images, structured values
- TTL support
- Resource associations

### Type Definition
Content type definitions with renderers, editors, and constraints.

### Product
Business products associated with tenants.

### Tool
External tool/integration configurations.

## MCP Tools

The MCP server exposes the following tools:

### Tenant Tools
- `create_tenant`: Create a new tenant
- `get_tenant`: Retrieve tenant information

### Library Tools
- `create_library`: Create a library in a tenant

### Notebook Tools
- `create_notebook`: Create a notebook in a library
- `append_content_block`: Add content blocks to notebooks

### Search Tools
- `semantic_search_notebooks`: Vector-based semantic search
- `graph_query_resources`: Graph-based relationship queries

## REST API Endpoints

### Tenants
- `GET /api/v1/tenants` - List all tenants
- `POST /api/v1/tenants` - Create a tenant
- `GET /api/v1/tenants/:id` - Get tenant by ID
- `PUT /api/v1/tenants/:id` - Update tenant
- `DELETE /api/v1/tenants/:id` - Delete tenant

### Libraries
- `GET /api/v1/libraries/by-tenant/:tenantId` - List libraries
- `POST /api/v1/libraries/by-tenant/:tenantId` - Create library
- `GET /api/v1/libraries/:id` - Get library
- `PUT /api/v1/libraries/:id` - Update library
- `DELETE /api/v1/libraries/:id` - Delete library

### Notebooks
- `GET /api/v1/notebooks/by-library/:libraryId` - List notebooks
- `POST /api/v1/notebooks/by-library/:libraryId` - Create notebook
- `GET /api/v1/notebooks/:id` - Get notebook
- `PUT /api/v1/notebooks/:id` - Update notebook
- `DELETE /api/v1/notebooks/:id` - Delete notebook

### Health
- `GET /health` - Health check endpoint

## Database Schema

The database includes:
- Core relational tables for all domain entities
- Vector embedding tables for semantic search
- Graph structure using Apache AGE for relationship queries
- Resource index for fast URI resolution
- Automatic triggers for graph synchronization

See `db/migrations/` for complete schema definitions.

## Development

### Running Tests
```bash
go test ./...
```

### Database Migrations

Migrations are automatically applied on startup. Migration files are located in `db/migrations/`:
- `001_create_schema.sql` - Core relational schema
- `002_create_graph_schema.sql` - Apache AGE graph setup
- `003_seed_data.sql` - Sample data for testing

### Adding a New Tool

1. Create tool definition in `internal/mcp/tools_*.go`:
   ```go
   func myNewTool() mcp.Tool {
       return mcp.Tool{
           Name: "my_new_tool",
           Description: "Description of the tool",
           InputSchema: mcp.ToolInputSchema{
               // Schema definition
           },
       }
   }
   ```

2. Implement handler:
   ```go
   func (s *Server) handleMyNewTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
       // Implementation
   }
   ```

3. Register in `server.go`:
   ```go
   s.mcpServer.AddTool(myNewTool(), s.handleMyNewTool)
   ```

## Roadmap

- [ ] Complete OIDC authentication integration
- [ ] Complete OPA authorization policies
- [ ] WebSocket support for real-time updates
- [ ] Embedding service integration for semantic search
- [ ] Enhanced graph query capabilities
- [ ] MCP resource templates
- [ ] Streamable HTTP transport
- [ ] Plugin architecture for external MCP servers
- [ ] Comprehensive test suite
- [ ] Performance optimization

## License

See LICENSE file for details.

## Contributing

Contributions are welcome! Please see CONTRIBUTING.md for guidelines.

## Support

For issues and questions, please use the GitHub issue tracker.

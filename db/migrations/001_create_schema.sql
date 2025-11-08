-- Migration 001: Create base schema for Synthesis
-- This creates all core tables for tenants, libraries, notebooks, features, types, products, and tools

-- Enable required extensions (should be done by init script, but ensure here too)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS vector;

-- Tenant table
CREATE TABLE IF NOT EXISTS tenant (
    id VARCHAR(255) PRIMARY KEY,
    owner VARCHAR(255) NOT NULL,
    display_name VARCHAR(500) NOT NULL,
    description TEXT,
    labels_json JSONB DEFAULT '{}',
    version VARCHAR(100),
    last_modified TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tenant_owner ON tenant(owner);
CREATE INDEX idx_tenant_labels ON tenant USING gin(labels_json);

-- Library table
CREATE TABLE IF NOT EXISTS library (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    owner VARCHAR(255) NOT NULL,
    display_name VARCHAR(500) NOT NULL,
    description TEXT,
    labels_json JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tenant_id, id)
);

CREATE INDEX idx_library_tenant ON library(tenant_id);
CREATE INDEX idx_library_owner ON library(owner);
CREATE INDEX idx_library_labels ON library USING gin(labels_json);

-- Notebook table
CREATE TABLE IF NOT EXISTS notebook (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    library_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'draft',
    owner VARCHAR(255) NOT NULL,
    display_name VARCHAR(500) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (library_id, tenant_id) REFERENCES library(id, tenant_id) ON DELETE CASCADE
);

CREATE INDEX idx_notebook_tenant ON notebook(tenant_id);
CREATE INDEX idx_notebook_library ON notebook(library_id);
CREATE INDEX idx_notebook_owner ON notebook(owner);
CREATE INDEX idx_notebook_status ON notebook(status);

-- Notebook content table
CREATE TABLE IF NOT EXISTS notebook_content (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    notebook_id VARCHAR(255) NOT NULL REFERENCES notebook(id) ON DELETE CASCADE,
    markdown TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notebook_content_notebook ON notebook_content(notebook_id);

-- Content block table
CREATE TABLE IF NOT EXISTS content_block (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    notebook_id VARCHAR(255) NOT NULL REFERENCES notebook(id) ON DELETE CASCADE,
    uid VARCHAR(255) NOT NULL,
    parent_uid VARCHAR(255),
    content_type VARCHAR(100) NOT NULL,
    data TEXT,
    "order" INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(notebook_id, uid)
);

CREATE INDEX idx_content_block_notebook ON content_block(notebook_id);
CREATE INDEX idx_content_block_parent ON content_block(parent_uid);
CREATE INDEX idx_content_block_order ON content_block("order");

-- Content block type mapping table
CREATE TABLE IF NOT EXISTS content_block_type (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content_block_id UUID NOT NULL REFERENCES content_block(id) ON DELETE CASCADE,
    type_name VARCHAR(255) NOT NULL
);

CREATE INDEX idx_content_block_type_block ON content_block_type(content_block_id);
CREATE INDEX idx_content_block_type_name ON content_block_type(type_name);

-- Notebook notification table
CREATE TABLE IF NOT EXISTS notebook_notification (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    notebook_id VARCHAR(255) NOT NULL REFERENCES notebook(id) ON DELETE CASCADE,
    nurl TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notebook_notification_notebook ON notebook_notification(notebook_id);

-- Feature table
CREATE TABLE IF NOT EXISTS feature (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    display_name VARCHAR(500) NOT NULL,
    description TEXT,
    ttl INTERVAL,
    values_json JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_feature_tenant ON feature(tenant_id);
CREATE INDEX idx_feature_expires ON feature(expires_at) WHERE expires_at IS NOT NULL;
CREATE INDEX idx_feature_values ON feature USING gin(values_json);

-- Feature resource table
CREATE TABLE IF NOT EXISTS feature_resource (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    feature_id VARCHAR(255) NOT NULL REFERENCES feature(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_feature_resource_feature ON feature_resource(feature_id);

-- Feature notification table
CREATE TABLE IF NOT EXISTS feature_notification (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    feature_id VARCHAR(255) NOT NULL REFERENCES feature(id) ON DELETE CASCADE,
    nurl TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_feature_notification_feature ON feature_notification(feature_id);

-- Type definition table
CREATE TABLE IF NOT EXISTS type_def (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    renderers_json JSONB DEFAULT '[]',
    editors_json JSONB DEFAULT '[]',
    constraints_json JSONB DEFAULT '[]',
    labels_json JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_type_def_name ON type_def(name);
CREATE INDEX idx_type_def_labels ON type_def USING gin(labels_json);

-- Product table
CREATE TABLE IF NOT EXISTS product (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    display_name VARCHAR(500) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_product_tenant ON product(tenant_id);

-- Product user table
CREATE TABLE IF NOT EXISTS product_user (
    product_id VARCHAR(255) NOT NULL REFERENCES product(id) ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL,
    role VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (product_id, user_id)
);

CREATE INDEX idx_product_user_user ON product_user(user_id);

-- Tool configuration table
CREATE TABLE IF NOT EXISTS tool (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    display_name VARCHAR(500) NOT NULL,
    description TEXT,
    config_json JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tool_tenant ON tool(tenant_id);

-- Resource index table for fast URI resolution
CREATE TABLE IF NOT EXISTS resource_index (
    uri TEXT PRIMARY KEY,
    entity_type VARCHAR(100) NOT NULL,
    entity_id VARCHAR(255) NOT NULL,
    tenant_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_resource_index_entity ON resource_index(entity_type, entity_id);
CREATE INDEX idx_resource_index_tenant ON resource_index(tenant_id);

-- Notebook embedding table for vector search
CREATE TABLE IF NOT EXISTS notebook_embedding (
    notebook_id VARCHAR(255) PRIMARY KEY REFERENCES notebook(id) ON DELETE CASCADE,
    embedding vector(1536), -- OpenAI ada-002 dimension, adjust as needed
    model VARCHAR(100) DEFAULT 'text-embedding-ada-002',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create vector index for similarity search
CREATE INDEX idx_notebook_embedding_vector ON notebook_embedding USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);

-- Feature embedding table for vector search
CREATE TABLE IF NOT EXISTS feature_embedding (
    feature_id VARCHAR(255) PRIMARY KEY REFERENCES feature(id) ON DELETE CASCADE,
    embedding vector(1536),
    model VARCHAR(100) DEFAULT 'text-embedding-ada-002',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_feature_embedding_vector ON feature_embedding USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply updated_at triggers to relevant tables
CREATE TRIGGER update_tenant_updated_at BEFORE UPDATE ON tenant
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_notebook_updated_at BEFORE UPDATE ON notebook
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_notebook_content_updated_at BEFORE UPDATE ON notebook_content
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_feature_updated_at BEFORE UPDATE ON feature
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_type_def_updated_at BEFORE UPDATE ON type_def
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_product_updated_at BEFORE UPDATE ON product
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tool_updated_at BEFORE UPDATE ON tool
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_resource_index_updated_at BEFORE UPDATE ON resource_index
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_notebook_embedding_updated_at BEFORE UPDATE ON notebook_embedding
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_feature_embedding_updated_at BEFORE UPDATE ON feature_embedding
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

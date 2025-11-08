-- Migration 002: Create graph schema using Apache AGE
-- This sets up the graph database for relationship queries

-- Ensure AGE extension is loaded
LOAD 'age';
SET search_path = ag_catalog, "$user", public;

-- Create the synthesis graph (if not already created by init script)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM ag_catalog.ag_graph WHERE name = 'synthesis_graph'
    ) THEN
        PERFORM ag_catalog.create_graph('synthesis_graph');
    END IF;
END
$$;

-- Create vertex labels (node types)
SELECT * FROM ag_catalog.create_vlabel('synthesis_graph', 'Tenant');
SELECT * FROM ag_catalog.create_vlabel('synthesis_graph', 'Library');
SELECT * FROM ag_catalog.create_vlabel('synthesis_graph', 'Notebook');
SELECT * FROM ag_catalog.create_vlabel('synthesis_graph', 'Feature');
SELECT * FROM ag_catalog.create_vlabel('synthesis_graph', 'Product');
SELECT * FROM ag_catalog.create_vlabel('synthesis_graph', 'Tool');
SELECT * FROM ag_catalog.create_vlabel('synthesis_graph', 'Type');

-- Create edge labels (relationship types)
SELECT * FROM ag_catalog.create_elabel('synthesis_graph', 'BELONGS_TO');
SELECT * FROM ag_catalog.create_elabel('synthesis_graph', 'DERIVES_FROM');
SELECT * FROM ag_catalog.create_elabel('synthesis_graph', 'PART_OF_PRODUCT');
SELECT * FROM ag_catalog.create_elabel('synthesis_graph', 'USES_TOOL');
SELECT * FROM ag_catalog.create_elabel('synthesis_graph', 'HAS_TYPE');
SELECT * FROM ag_catalog.create_elabel('synthesis_graph', 'RELATED_TO');

-- Create helper functions for graph operations

-- Function to sync tenant to graph
CREATE OR REPLACE FUNCTION sync_tenant_to_graph()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT' OR TG_OP = 'UPDATE') THEN
        -- Upsert vertex in graph
        PERFORM ag_catalog.cypher(
            'synthesis_graph',
            $$MERGE (t:Tenant {id: $tenant_id})
              SET t.display_name = $display_name,
                  t.owner = $owner,
                  t.description = $description$$,
            agtype_build_map(
                'tenant_id', NEW.id::agtype,
                'display_name', NEW.display_name::agtype,
                'owner', NEW.owner::agtype,
                'description', COALESCE(NEW.description, '')::agtype
            )
        );
    ELSIF (TG_OP = 'DELETE') THEN
        -- Delete vertex from graph
        PERFORM ag_catalog.cypher(
            'synthesis_graph',
            $$MATCH (t:Tenant {id: $tenant_id})
              DETACH DELETE t$$,
            agtype_build_map('tenant_id', OLD.id::agtype)
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to sync library to graph
CREATE OR REPLACE FUNCTION sync_library_to_graph()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT' OR TG_OP = 'UPDATE') THEN
        -- Upsert library vertex and relationship
        PERFORM ag_catalog.cypher(
            'synthesis_graph',
            $$MERGE (l:Library {id: $library_id})
              SET l.display_name = $display_name,
                  l.owner = $owner
              WITH l
              MATCH (t:Tenant {id: $tenant_id})
              MERGE (l)-[:BELONGS_TO]->(t)$$,
            agtype_build_map(
                'library_id', NEW.id::agtype,
                'tenant_id', NEW.tenant_id::agtype,
                'display_name', NEW.display_name::agtype,
                'owner', NEW.owner::agtype
            )
        );
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM ag_catalog.cypher(
            'synthesis_graph',
            $$MATCH (l:Library {id: $library_id})
              DETACH DELETE l$$,
            agtype_build_map('library_id', OLD.id::agtype)
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to sync notebook to graph
CREATE OR REPLACE FUNCTION sync_notebook_to_graph()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT' OR TG_OP = 'UPDATE') THEN
        PERFORM ag_catalog.cypher(
            'synthesis_graph',
            $$MERGE (n:Notebook {id: $notebook_id})
              SET n.display_name = $display_name,
                  n.owner = $owner,
                  n.status = $status
              WITH n
              MATCH (l:Library {id: $library_id})
              MERGE (n)-[:BELONGS_TO]->(l)$$,
            agtype_build_map(
                'notebook_id', NEW.id::agtype,
                'library_id', NEW.library_id::agtype,
                'display_name', NEW.display_name::agtype,
                'owner', NEW.owner::agtype,
                'status', NEW.status::agtype
            )
        );
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM ag_catalog.cypher(
            'synthesis_graph',
            $$MATCH (n:Notebook {id: $notebook_id})
              DETACH DELETE n$$,
            agtype_build_map('notebook_id', OLD.id::agtype)
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers to keep graph in sync
CREATE TRIGGER sync_tenant_to_graph_trigger
    AFTER INSERT OR UPDATE OR DELETE ON tenant
    FOR EACH ROW EXECUTE FUNCTION sync_tenant_to_graph();

CREATE TRIGGER sync_library_to_graph_trigger
    AFTER INSERT OR UPDATE OR DELETE ON library
    FOR EACH ROW EXECUTE FUNCTION sync_library_to_graph();

CREATE TRIGGER sync_notebook_to_graph_trigger
    AFTER INSERT OR UPDATE OR DELETE ON notebook
    FOR EACH ROW EXECUTE FUNCTION sync_notebook_to_graph();

-- Similar functions can be created for feature, product, tool as needed
-- For now, we'll handle those in application code

-- Migration 003: Seed data for development and testing

-- Insert example tenant
INSERT INTO tenant (id, owner, display_name, description, labels_json, version)
VALUES (
    'technicaldetails',
    'josh@prismon.com',
    'Joshua''s Personal Tenant',
    'This is a simple tenant',
    '{"environment": "production", "app": "nginx"}',
    '1234567'
) ON CONFLICT (id) DO NOTHING;

-- Insert example library
INSERT INTO library (id, tenant_id, owner, display_name, description)
VALUES (
    'technicaldetails-default-library',
    'technicaldetails',
    'josh@prismon.com',
    'Default Library',
    'Default library for technical details'
) ON CONFLICT (id) DO NOTHING;

-- Insert example notebook
INSERT INTO notebook (id, tenant_id, library_id, status, owner, display_name, description)
VALUES (
    'josh_notes_2025',
    'technicaldetails',
    'technicaldetails-default-library',
    'approved',
    'josh@prismon.com',
    'Notes from 3/27/2025',
    'Here are knowledge elements from my obsidian notebook'
) ON CONFLICT (id) DO NOTHING;

-- Insert notebook content
INSERT INTO notebook_content (notebook_id, markdown)
SELECT 'josh_notes_2025', '## Typographic replacements

Enable typographer option to see result.
'
WHERE NOT EXISTS (
    SELECT 1 FROM notebook_content WHERE notebook_id = 'josh_notes_2025'
);

-- Insert content blocks
INSERT INTO content_block (notebook_id, uid, content_type, data, "order")
VALUES
    (
        'josh_notes_2025',
        '550e8400-e29b-41d4-a716-446655440000',
        'text/markdown',
        '# Header
This is the first paragraph',
        1
    ),
    (
        'josh_notes_2025',
        '6ba7b810-9dad-11d1-80b4-00c04fd430c8',
        'binary/image',
        'base64encodedimagedata',
        2
    )
ON CONFLICT (notebook_id, uid) DO NOTHING;

-- Link content block to parent
UPDATE content_block
SET parent_uid = '550e8400-e29b-41d4-a716-446655440000'
WHERE uid = '6ba7b810-9dad-11d1-80b4-00c04fd430c8';

-- Insert content block types
INSERT INTO content_block_type (content_block_id, type_name)
SELECT cb.id, 'heading'
FROM content_block cb
WHERE cb.uid = '550e8400-e29b-41d4-a716-446655440000'
AND NOT EXISTS (
    SELECT 1 FROM content_block_type cbt
    WHERE cbt.content_block_id = cb.id AND cbt.type_name = 'heading'
);

INSERT INTO content_block_type (content_block_id, type_name)
SELECT cb.id, 'text'
FROM content_block cb
WHERE cb.uid = '550e8400-e29b-41d4-a716-446655440000'
AND NOT EXISTS (
    SELECT 1 FROM content_block_type cbt
    WHERE cbt.content_block_id = cb.id AND cbt.type_name = 'text'
);

INSERT INTO content_block_type (content_block_id, type_name)
SELECT cb.id, 'image'
FROM content_block cb
WHERE cb.uid = '6ba7b810-9dad-11d1-80b4-00c04fd430c8'
AND NOT EXISTS (
    SELECT 1 FROM content_block_type cbt
    WHERE cbt.content_block_id = cb.id AND cbt.type_name = 'image'
);

-- Insert notebook notification
INSERT INTO notebook_notification (notebook_id, nurl)
VALUES ('josh_notes_2025', 'https://www.technicaldetails.org/modify_alert')
ON CONFLICT DO NOTHING;

-- Insert example type definitions
INSERT INTO type_def (name, description, renderers_json, editors_json, constraints_json)
VALUES
    (
        'markdown',
        'Markdown content type',
        '[{"name": "markdown-renderer", "config": {}}]',
        '[{"name": "markdown-editor", "config": {"toolbar": true}}]',
        '[{"type": "max-length", "config": {"max": 100000}}]'
    ),
    (
        'image',
        'Image content type',
        '[{"name": "image-renderer", "config": {"maxWidth": 800}}]',
        '[{"name": "image-uploader", "config": {"formats": ["jpg", "png", "gif", "webp"]}}]',
        '[{"type": "file-size", "config": {"maxBytes": 5242880}}]'
    ),
    (
        'json',
        'JSON data type',
        '[{"name": "json-viewer", "config": {"pretty": true}}]',
        '[{"name": "json-editor", "config": {"schema": null}}]',
        '[{"type": "json-schema", "config": {}}]'
    )
ON CONFLICT (name) DO NOTHING;

-- Insert example product
INSERT INTO product (id, tenant_id, display_name, description)
VALUES (
    'product-alpha',
    'technicaldetails',
    'Product Alpha',
    'First product for technical details tenant'
) ON CONFLICT (id) DO NOTHING;

-- Insert product users
INSERT INTO product_user (product_id, user_id, role)
VALUES
    ('product-alpha', 'josh@prismon.com', 'owner'),
    ('product-alpha', 'user2@example.com', 'contributor')
ON CONFLICT (product_id, user_id) DO NOTHING;

-- Insert example tool configuration
INSERT INTO tool (id, tenant_id, display_name, description, config_json)
VALUES (
    'slack-integration',
    'technicaldetails',
    'Slack Integration',
    'Integration with Slack workspace',
    '{"webhook_url": "https://hooks.slack.com/services/xxx", "channel": "#general"}'
) ON CONFLICT (id) DO NOTHING;

-- Insert resource index entries
INSERT INTO resource_index (uri, entity_type, entity_id, tenant_id)
VALUES
    ('synthesis://tenant/technicaldetails', 'tenant', 'technicaldetails', 'technicaldetails'),
    ('synthesis://tenant/technicaldetails/library/technicaldetails-default-library', 'library', 'technicaldetails-default-library', 'technicaldetails'),
    ('synthesis://tenant/technicaldetails/notebook/josh_notes_2025', 'notebook', 'josh_notes_2025', 'technicaldetails'),
    ('synthesis://tenant/technicaldetails/product/product-alpha', 'product', 'product-alpha', 'technicaldetails'),
    ('synthesis://tenant/technicaldetails/tool/slack-integration', 'tool', 'slack-integration', 'technicaldetails'),
    ('synthesis://type/markdown', 'type', 'markdown', NULL),
    ('synthesis://type/image', 'type', 'image', NULL),
    ('synthesis://type/json', 'type', 'json', NULL)
ON CONFLICT (uri) DO NOTHING;

package domain

// Notebook represents the primary editable element with hierarchical content blocks
type Notebook struct {
	TenantID      string              `json:"tenantId" db:"tenant_id"`
	ID            string              `json:"notebookId" db:"id"`
	LibraryID     string              `json:"libraryId" db:"library_id"`
	Status        string              `json:"status" db:"status"`
	Owner         string              `json:"owner" db:"owner"`
	DisplayName   string              `json:"display_name" db:"display_name"`
	Description   string              `json:"description" db:"description"`
	Contents      NotebookContents    `json:"contents"`
	Notifications []Notification      `json:"notification,omitempty"`
}

// NotebookContents represents the content structure of a notebook
type NotebookContents struct {
	Data          NotebookData   `json:"data"`
	ContentBlocks []ContentBlock `json:"content_blocks,omitempty"`
}

// NotebookData contains the main markdown content
type NotebookData struct {
	Markdown string `json:"Markdown" db:"markdown"`
}

// ContentBlock represents a hierarchical content block within a notebook
type ContentBlock struct {
	UID         string   `json:"uid" db:"uid"`
	ParentUID   *string  `json:"parent_uid,omitempty" db:"parent_uid"`
	ContentType string   `json:"content_type" db:"content_type"`
	Data        string   `json:"data" db:"data"`
	Order       int      `json:"order" db:"order"`
	Types       []string `json:"types" db:"-"`
}

// Notification represents a webhook URL for notifications
type Notification struct {
	URL string `json:"nurl" db:"nurl"`
}

// URI returns the MCP URI for this notebook
func (n *Notebook) URI() string {
	return "synthesis://tenant/" + n.TenantID + "/notebook/" + n.ID
}

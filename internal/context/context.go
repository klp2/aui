package context

// Context represents a collection of files and metadata for AI agent consumption
type Context struct {
	Name        string
	Description string
	Files       []*File
	TotalTokens int
}

// NewContext creates a new context with the given name and description
func NewContext(name, description string) *Context {
	return &Context{
		Name:        name,
		Description: description,
		Files:       []*File{},
		TotalTokens: 0,
	}
}

// AddFile adds a file to the context, replacing if it already exists with different hash
func (c *Context) AddFile(file *File) {
	if file == nil {
		return
	}

	// Check if file with same path already exists
	for i, f := range c.Files {
		if f.Path == file.Path {
			if f.Equals(file) {
				// Same file (same path and hash), don't add duplicate
				return
			}
			// Different hash, replace the old version
			c.TotalTokens -= f.TokenCount
			c.Files[i] = file
			c.TotalTokens += file.TokenCount
			return
		}
	}

	// New file, add it
	c.Files = append(c.Files, file)
	c.TotalTokens += file.TokenCount
}

// RemoveFile removes a file from the context by path
func (c *Context) RemoveFile(path string) {
	newFiles := []*File{}
	for _, f := range c.Files {
		if f.Path != path {
			newFiles = append(newFiles, f)
		} else {
			// Subtract tokens when removing
			c.TotalTokens -= f.TokenCount
		}
	}
	c.Files = newFiles
}

// GetFile returns a file by path, or nil if not found
func (c *Context) GetFile(path string) *File {
	for _, f := range c.Files {
		if f.Path == path {
			return f
		}
	}
	return nil
}

// HasFile checks if a file with the given path exists in the context
func (c *Context) HasFile(path string) bool {
	return c.GetFile(path) != nil
}

// Clear removes all files and resets token count, but keeps name and description
func (c *Context) Clear() {
	c.Files = []*File{}
	c.TotalTokens = 0
}

// RecalculateTokens recalculates the total token count from all files
func (c *Context) RecalculateTokens() {
	total := 0
	for _, f := range c.Files {
		total += f.TokenCount
	}
	c.TotalTokens = total
}

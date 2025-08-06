package context

import (
	"crypto/rand"
	"encoding/hex"
	"path/filepath"
	"strings"
	"time"
)

// File represents a file in a context with its metadata
type File struct {
	ID         string
	Path       string
	Name       string
	Content    string
	Size       int64
	Hash       string
	Language   string
	Tokens     int       // Renamed from TokenCount for consistency with storage
	ModifiedAt time.Time // Renamed from LastModified for consistency
}

// NewFile creates a new file with the given path and name
func NewFile(path, name string) *File {
	return &File{
		ID:         generateFileID(),
		Path:       path,
		Name:       name,
		ModifiedAt: time.Now(),
	}
}

// generateFileID generates a random ID for a file
func generateFileID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// UpdateMetadata updates all the file's metadata
func (f *File) UpdateMetadata(size int64, hash, language string, tokens int, modifiedAt time.Time) {
	f.Size = size
	f.Hash = hash
	f.Language = language
	f.Tokens = tokens
	f.ModifiedAt = modifiedAt
}

// DetectLanguage detects the programming language based on file extension
func (f *File) DetectLanguage() {
	ext := strings.ToLower(filepath.Ext(f.Path))
	base := filepath.Base(f.Path)

	// Check special filenames first
	switch base {
	case "Makefile", "makefile":
		f.Language = "makefile"
		return
	}

	// Check by extension
	languageMap := map[string]string{
		".go":   "go",
		".py":   "python",
		".js":   "javascript",
		".jsx":  "javascript",
		".ts":   "typescript",
		".tsx":  "typescript",
		".css":  "css",
		".toml": "toml",
		".md":   "markdown",
		".rs":   "rust",
		".c":    "c",
		".cpp":  "cpp",
		".h":    "c",
		".hpp":  "cpp",
		".java": "java",
		".rb":   "ruby",
		".php":  "php",
		".sh":   "shell",
		".bash": "shell",
		".zsh":  "shell",
		".yaml": "yaml",
		".yml":  "yaml",
		".json": "json",
		".xml":  "xml",
		".html": "html",
		".htm":  "html",
	}

	if lang, ok := languageMap[ext]; ok {
		f.Language = lang
	} else {
		f.Language = ""
	}
}

// Equals checks if two files are the same (same path and hash)
func (f *File) Equals(other *File) bool {
	if other == nil {
		return false
	}
	return f.Path == other.Path && f.Hash == other.Hash
}

// NeedsUpdate checks if the file needs to be updated based on hash or modification time
func (f *File) NeedsUpdate(newHash string, newModTime time.Time) bool {
	// If hash is different, needs update
	if f.Hash != newHash {
		return true
	}

	// If modification time is newer, needs update
	if newModTime.After(f.ModifiedAt) {
		return true
	}

	return false
}

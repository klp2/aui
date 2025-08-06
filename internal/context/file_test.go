package context

import (
	"testing"
	"time"
)

func TestNewFile(t *testing.T) {
	path := "auth/login.go"
	name := "login.go"

	file := NewFile(path, name)

	if file.ID == "" {
		t.Error("NewFile().ID should not be empty")
	}

	if file.Path != path {
		t.Errorf("NewFile().Path = %v, want %v", file.Path, path)
	}

	if file.Name != name {
		t.Errorf("NewFile().Name = %v, want %v", file.Name, name)
	}

	if file.Size != 0 {
		t.Errorf("NewFile().Size = %v, want 0", file.Size)
	}

	if file.Hash != "" {
		t.Errorf("NewFile().Hash = %v, want empty", file.Hash)
	}

	if file.Language != "" {
		t.Errorf("NewFile().Language = %v, want empty", file.Language)
	}

	if file.Tokens != 0 {
		t.Errorf("NewFile().Tokens = %v, want 0", file.Tokens)
	}

	if file.ModifiedAt.IsZero() {
		t.Errorf("NewFile().ModifiedAt should not be zero")
	}
}

func TestFileUpdateMetadata(t *testing.T) {
	file := NewFile("main.go", "main.go")

	size := int64(1024)
	hash := "abc123def456"
	lang := "go"
	tokens := 150
	modTime := time.Now().Add(-1 * time.Hour)

	file.UpdateMetadata(size, hash, lang, tokens, modTime)

	if file.Size != size {
		t.Errorf("After UpdateMetadata(), Size = %v, want %v", file.Size, size)
	}

	if file.Hash != hash {
		t.Errorf("After UpdateMetadata(), Hash = %v, want %v", file.Hash, hash)
	}

	if file.Language != lang {
		t.Errorf("After UpdateMetadata(), Language = %v, want %v", file.Language, lang)
	}

	if file.Tokens != tokens {
		t.Errorf("After UpdateMetadata(), Tokens = %v, want %v", file.Tokens, tokens)
	}

	if !file.ModifiedAt.Equal(modTime) {
		t.Errorf("After UpdateMetadata(), ModifiedAt = %v, want %v", file.ModifiedAt, modTime)
	}
}

func TestFileDetectLanguage(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantLang string
	}{
		{"Go file", "main.go", "go"},
		{"Python file", "script.py", "python"},
		{"JavaScript file", "app.js", "javascript"},
		{"TypeScript file", "app.ts", "typescript"},
		{"Rust file", "main.rs", "rust"},
		{"C file", "program.c", "c"},
		{"C++ file", "program.cpp", "cpp"},
		{"Java file", "Main.java", "java"},
		{"Ruby file", "script.rb", "ruby"},
		{"Shell script", "deploy.sh", "shell"},
		{"Makefile", "Makefile", "makefile"},
		{"Lowercase makefile", "makefile", "makefile"},
		{"JSON file", "config.json", "json"},
		{"YAML file", "config.yaml", "yaml"},
		{"Markdown file", "README.md", "markdown"},
		{"Unknown extension", "data.xyz", ""},
		{"No extension", "README", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := NewFile(tt.path, tt.path)
			file.DetectLanguage()

			if file.Language != tt.wantLang {
				t.Errorf("DetectLanguage() for %v = %v, want %v", tt.path, file.Language, tt.wantLang)
			}
		})
	}
}

func TestFileEquals(t *testing.T) {
	// Test files with same path and hash
	file1 := NewFile("main.go", "main.go")
	file1.Hash = "abc123"

	file2 := NewFile("main.go", "main.go")
	file2.Hash = "abc123"

	if !file1.Equals(file2) {
		t.Error("Files with same path and hash should be equal")
	}

	// Test files with same path, different hash
	file3 := NewFile("main.go", "main.go")
	file3.Hash = "def456"

	if file1.Equals(file3) {
		t.Error("Files with same path but different hash should not be equal")
	}

	// Test files with different path, same hash
	file4 := NewFile("other.go", "other.go")
	file4.Hash = "abc123"

	if file1.Equals(file4) {
		t.Error("Files with different path should not be equal")
	}

	// Test nil file
	if file1.Equals(nil) {
		t.Error("File should not equal nil")
	}
}

func TestFileNeedsUpdate(t *testing.T) {
	baseTime := time.Now()
	file := NewFile("main.go", "main.go")
	file.Hash = "abc123"
	file.ModifiedAt = baseTime

	tests := []struct {
		name       string
		newHash    string
		newModTime time.Time
		wantUpdate bool
	}{
		{
			name:       "Same hash, same time",
			newHash:    "abc123",
			newModTime: baseTime,
			wantUpdate: false,
		},
		{
			name:       "Different hash, same time",
			newHash:    "def456",
			newModTime: baseTime,
			wantUpdate: true,
		},
		{
			name:       "Same hash, newer time",
			newHash:    "abc123",
			newModTime: baseTime.Add(1 * time.Hour),
			wantUpdate: true,
		},
		{
			name:       "Same hash, older time",
			newHash:    "abc123",
			newModTime: baseTime.Add(-1 * time.Hour),
			wantUpdate: false,
		},
		{
			name:       "Different hash, newer time",
			newHash:    "def456",
			newModTime: baseTime.Add(1 * time.Hour),
			wantUpdate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := file.NeedsUpdate(tt.newHash, tt.newModTime)
			if result != tt.wantUpdate {
				t.Errorf("NeedsUpdate() = %v, want %v", result, tt.wantUpdate)
			}
		})
	}
}

func TestGenerateFileID(t *testing.T) {
	// Test that IDs are unique
	id1 := generateFileID()
	id2 := generateFileID()

	if id1 == id2 {
		t.Error("generateFileID() should produce unique IDs")
	}

	if len(id1) != 16 {
		t.Errorf("generateFileID() should produce 16-character hex strings, got %d", len(id1))
	}
}

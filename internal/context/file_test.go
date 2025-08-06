package context

import (
	"testing"
	"time"
)

func TestNewFile(t *testing.T) {
	path := "auth/login.go"

	file := NewFile(path)

	if file.Path != path {
		t.Errorf("NewFile().Path = %v, want %v", file.Path, path)
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

	if file.TokenCount != 0 {
		t.Errorf("NewFile().TokenCount = %v, want 0", file.TokenCount)
	}

	if file.LastModified.IsZero() {
		t.Errorf("NewFile().LastModified should not be zero")
	}
}

func TestFileUpdateMetadata(t *testing.T) {
	file := NewFile("main.go")

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

	if file.TokenCount != tokens {
		t.Errorf("After UpdateMetadata(), TokenCount = %v, want %v", file.TokenCount, tokens)
	}

	if !file.LastModified.Equal(modTime) {
		t.Errorf("After UpdateMetadata(), LastModified = %v, want %v", file.LastModified, modTime)
	}
}

func TestFileDetectLanguage(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"main.go", "go"},
		{"app.py", "python"},
		{"index.js", "javascript"},
		{"App.tsx", "typescript"},
		{"style.css", "css"},
		{"Cargo.toml", "toml"},
		{"README.md", "markdown"},
		{"Makefile", "makefile"},
		{"unknown.xyz", ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			file := NewFile(tt.path)
			file.DetectLanguage()

			if file.Language != tt.expected {
				t.Errorf("DetectLanguage() for %v = %v, want %v", tt.path, file.Language, tt.expected)
			}
		})
	}
}

func TestFileEquals(t *testing.T) {
	file1 := NewFile("auth/login.go")
	file1.UpdateMetadata(1024, "hash1", "go", 100, time.Now())

	file2 := NewFile("auth/login.go")
	file2.UpdateMetadata(1024, "hash1", "go", 100, time.Now())

	file3 := NewFile("auth/session.go")
	file3.UpdateMetadata(1024, "hash1", "go", 100, time.Now())

	file4 := NewFile("auth/login.go")
	file4.UpdateMetadata(1024, "hash2", "go", 100, time.Now())

	tests := []struct {
		name     string
		file1    *File
		file2    *File
		expected bool
	}{
		{"same path and hash", file1, file2, true},
		{"different path", file1, file3, false},
		{"same path different hash", file1, file4, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.file1.Equals(tt.file2)
			if result != tt.expected {
				t.Errorf("Equals() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFileNeedsUpdate(t *testing.T) {
	oldTime := time.Now().Add(-2 * time.Hour)
	newTime := time.Now()

	file := NewFile("main.go")
	file.UpdateMetadata(1024, "hash1", "go", 100, oldTime)

	// File with same hash shouldn't need update
	if file.NeedsUpdate("hash1", oldTime) {
		t.Error("NeedsUpdate() = true for same hash and time, want false")
	}

	// File with different hash should need update
	if !file.NeedsUpdate("hash2", oldTime) {
		t.Error("NeedsUpdate() = false for different hash, want true")
	}

	// File with newer modification time should need update
	if !file.NeedsUpdate("hash1", newTime) {
		t.Error("NeedsUpdate() = false for newer modification time, want true")
	}
}

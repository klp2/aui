package context

import (
	"testing"
	"time"
)

func TestNewContext(t *testing.T) {
	name := "bug-fix-auth"
	description := "Authentication bug context"

	ctx := NewContext(name, description)

	if ctx.Name != name {
		t.Errorf("NewContext().Name = %v, want %v", ctx.Name, name)
	}

	if ctx.Description != description {
		t.Errorf("NewContext().Description = %v, want %v", ctx.Description, description)
	}

	if len(ctx.Files) != 0 {
		t.Errorf("NewContext().Files length = %v, want 0", len(ctx.Files))
	}

	if ctx.TotalTokens != 0 {
		t.Errorf("NewContext().TotalTokens = %v, want 0", ctx.TotalTokens)
	}
}

func TestContextAddFile(t *testing.T) {
	ctx := NewContext("test", "test context")

	file1 := NewFile("auth/login.go")
	file1.UpdateMetadata(1024, "hash1", "go", 100, time.Now())

	file2 := NewFile("auth/session.go")
	file2.UpdateMetadata(2048, "hash2", "go", 200, time.Now())

	ctx.AddFile(file1)
	ctx.AddFile(file2)

	if len(ctx.Files) != 2 {
		t.Errorf("After adding 2 files, Files length = %v, want 2", len(ctx.Files))
	}

	if ctx.Files[0].Path != file1.Path {
		t.Errorf("First file path = %v, want %v", ctx.Files[0].Path, file1.Path)
	}

	if ctx.Files[1].Path != file2.Path {
		t.Errorf("Second file path = %v, want %v", ctx.Files[1].Path, file2.Path)
	}

	// Check total tokens are calculated
	expectedTokens := 100 + 200
	if ctx.TotalTokens != expectedTokens {
		t.Errorf("After adding files, TotalTokens = %v, want %v", ctx.TotalTokens, expectedTokens)
	}
}

func TestContextAddFileDuplicate(t *testing.T) {
	ctx := NewContext("test", "test context")

	file := NewFile("auth/login.go")
	file.UpdateMetadata(1024, "hash1", "go", 100, time.Now())

	ctx.AddFile(file)
	ctx.AddFile(file) // Try to add duplicate (same path and hash)

	if len(ctx.Files) != 1 {
		t.Errorf("After adding duplicate file, Files length = %v, want 1", len(ctx.Files))
	}

	if ctx.TotalTokens != 100 {
		t.Errorf("After adding duplicate, TotalTokens = %v, want 100", ctx.TotalTokens)
	}
}

func TestContextAddFileUpdated(t *testing.T) {
	ctx := NewContext("test", "test context")

	// Add original file
	file1 := NewFile("auth/login.go")
	file1.UpdateMetadata(1024, "hash1", "go", 100, time.Now())
	ctx.AddFile(file1)

	// Add updated version of same file (same path, different hash)
	file2 := NewFile("auth/login.go")
	file2.UpdateMetadata(2048, "hash2", "go", 150, time.Now().Add(1*time.Hour))
	ctx.AddFile(file2)

	// Should replace the old version
	if len(ctx.Files) != 1 {
		t.Errorf("After adding updated file, Files length = %v, want 1", len(ctx.Files))
	}

	if ctx.Files[0].Hash != "hash2" {
		t.Errorf("File hash = %v, want hash2", ctx.Files[0].Hash)
	}

	if ctx.TotalTokens != 150 {
		t.Errorf("After update, TotalTokens = %v, want 150", ctx.TotalTokens)
	}
}

func TestContextRemoveFile(t *testing.T) {
	ctx := NewContext("test", "test context")

	file1 := NewFile("auth/login.go")
	file1.UpdateMetadata(1024, "hash1", "go", 100, time.Now())

	file2 := NewFile("auth/session.go")
	file2.UpdateMetadata(2048, "hash2", "go", 200, time.Now())

	file3 := NewFile("auth/token.go")
	file3.UpdateMetadata(512, "hash3", "go", 50, time.Now())

	ctx.AddFile(file1)
	ctx.AddFile(file2)
	ctx.AddFile(file3)

	ctx.RemoveFile("auth/session.go")

	if len(ctx.Files) != 2 {
		t.Errorf("After removing file, Files length = %v, want 2", len(ctx.Files))
	}

	// Check that file2 was removed
	for _, f := range ctx.Files {
		if f.Path == "auth/session.go" {
			t.Errorf("File %v should have been removed", file2.Path)
		}
	}

	// Check token count updated
	expectedTokens := 100 + 50
	if ctx.TotalTokens != expectedTokens {
		t.Errorf("After removal, TotalTokens = %v, want %v", ctx.TotalTokens, expectedTokens)
	}
}

func TestContextClear(t *testing.T) {
	ctx := NewContext("test", "test context")

	file1 := NewFile("file1.go")
	file1.TokenCount = 100
	file2 := NewFile("file2.go")
	file2.TokenCount = 200

	ctx.AddFile(file1)
	ctx.AddFile(file2)

	ctx.Clear()

	if len(ctx.Files) != 0 {
		t.Errorf("After Clear(), Files length = %v, want 0", len(ctx.Files))
	}

	if ctx.TotalTokens != 0 {
		t.Errorf("After Clear(), TotalTokens = %v, want 0", ctx.TotalTokens)
	}

	// Name and Description should remain
	if ctx.Name != "test" {
		t.Errorf("After Clear(), Name = %v, want 'test'", ctx.Name)
	}

	if ctx.Description != "test context" {
		t.Errorf("After Clear(), Description = %v, want 'test context'", ctx.Description)
	}
}

func TestContextGetFile(t *testing.T) {
	ctx := NewContext("test", "test context")

	file1 := NewFile("auth/login.go")
	file1.UpdateMetadata(1024, "hash1", "go", 100, time.Now())

	file2 := NewFile("auth/session.go")
	file2.UpdateMetadata(2048, "hash2", "go", 200, time.Now())

	ctx.AddFile(file1)
	ctx.AddFile(file2)

	// Test getting existing file
	found := ctx.GetFile("auth/login.go")
	if found == nil {
		t.Error("GetFile() returned nil for existing file")
	} else if found.Path != "auth/login.go" {
		t.Errorf("GetFile() returned wrong file: %v", found.Path)
	}

	// Test getting non-existent file
	notFound := ctx.GetFile("auth/nonexistent.go")
	if notFound != nil {
		t.Errorf("GetFile() should return nil for non-existent file, got %v", notFound)
	}
}

func TestContextHasFile(t *testing.T) {
	ctx := NewContext("test", "test context")

	file := NewFile("auth/login.go")
	ctx.AddFile(file)

	if !ctx.HasFile("auth/login.go") {
		t.Error("HasFile() returned false for existing file")
	}

	if ctx.HasFile("auth/nonexistent.go") {
		t.Error("HasFile() returned true for non-existent file")
	}
}

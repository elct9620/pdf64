package builder_test

import (
	"testing"

	"github.com/elct9620/pdf64/internal/builder"
)

func TestFileBuilder_BuildFromPath(t *testing.T) {
	// Arrange
	fileBuilder := builder.NewFileBuilder()
	path := "/path/to/file.pdf"

	// Act
	file, err := fileBuilder.BuildFromPath(path)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if file == nil {
		t.Fatal("expected file to not be nil")
	}

	if file.Path() != path {
		t.Errorf("expected path to be %q, got %q", path, file.Path())
	}

	if file.Id() == "" {
		t.Error("expected ID to not be empty")
	}
}

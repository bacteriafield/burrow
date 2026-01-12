package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/elaurentium/burrow/lib/fs"
)

func TestNewCreator(t *testing.T) {
	creator := fs.NewCreator()

	if creator == nil {
		t.Fatal("NewCreator() returned nil")
	}

	if creator.Perm != 0755 {
		t.Errorf("Expected Perm to be 0755, got %o", creator.Perm)
	}

	if creator.Workers != 0 {
		t.Errorf("Expected Workers to be 0, got %d", creator.Workers)
	}

	if creator.Wg == nil {
		t.Error("Expected Wg to be initialized, got nil")
	}
}

func TestCreator_CreateFiles(t *testing.T) {
	tempDir := t.TempDir()

	creator := fs.NewCreator()

	tests := []struct {
		name    string
		paths   []string
		wantErr bool
	}{
		{
			name: "create single file",
			paths: []string{
				filepath.Join(tempDir, "test1.txt"),
			},
			wantErr: false,
		},
		{
			name: "create file with parent directories",
			paths: []string{
				filepath.Join(tempDir, "dir1", "dir2", "test2.txt"),
			},
			wantErr: false,
		},
		{
			name: "create multiple files",
			paths: []string{
				filepath.Join(tempDir, "file1.txt"),
				filepath.Join(tempDir, "file2.txt"),
				filepath.Join(tempDir, "file3.txt"),
			},
			wantErr: false,
		},
		{
			name: "create file in current directory",
			paths: []string{
				filepath.Join(tempDir, "simple.txt"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := creator.Create(tt.paths)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, path := range tt.paths {
				if _, err := os.Stat(path); os.IsNotExist(err) {
					t.Errorf("File %s was not created", path)
				}
			}
		})
	}
}

func TestCreator_CreateDirectories(t *testing.T) {
	tempDir := t.TempDir()

	creator := fs.NewCreator()

	tests := []struct {
		name    string
		paths   []string
		wantErr bool
	}{
		{
			name: "create single directory",
			paths: []string{
				filepath.Join(tempDir, "testdir"),
			},
			wantErr: false,
		},
		{
			name: "create nested directories",
			paths: []string{
				filepath.Join(tempDir, "dir1", "dir2", "dir3"),
			},
			wantErr: false,
		},
		{
			name: "create multiple directories",
			paths: []string{
				filepath.Join(tempDir, "dirA"),
				filepath.Join(tempDir, "dirB"),
				filepath.Join(tempDir, "dirC"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := creator.Create(tt.paths)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, path := range tt.paths {
				info, err := os.Stat(path)
				if os.IsNotExist(err) {
					t.Errorf("Directory %s was not created", path)
				}
				if err == nil && !info.IsDir() {
					t.Errorf("Path %s is not a directory", path)
				}
			}
		})
	}
}

func TestCreator_CreateMixed(t *testing.T) {
	tempDir := t.TempDir()

	creator := fs.NewCreator()

	paths := []string{
		filepath.Join(tempDir, "mixeddir/"),
		filepath.Join(tempDir, "mixeddir/", "file.txt"),
		filepath.Join(tempDir, "anotherdir/", "subdir"),
		filepath.Join(tempDir, "anotherdir/", "file2.txt"),
	}

	err := creator.Create(paths)
	if err != nil {
		t.Errorf("Create() unexpected error = %v", err)
	}

	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Path %s was not created", path)
		}
	}
}

func TestCreator_CreateDuplicateFile(t *testing.T) {
	tempDir := t.TempDir()

	creator := fs.NewCreator()

	filePath := filepath.Join(tempDir, "duplicate.txt")

	err := creator.Create([]string{filePath})
	if err != nil {
		t.Errorf("First Create() unexpected error = %v", err)
	}

	err = creator.Create([]string{filePath})
	if err != nil {
		t.Errorf("Second Create() unexpected error = %v", err)
	}
}

func TestCreator_CustomPermissions(t *testing.T) {
	tempDir := t.TempDir()

	creator := fs.NewCreator()
	creator.Perm = 0700

	filePath := filepath.Join(tempDir, "permtest.txt")
	dirPath := filepath.Join(tempDir, "permdir")

	err := creator.Create([]string{filePath, dirPath})
	if err != nil {
		t.Errorf("Create() unexpected error = %v", err)
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	if fileInfo.Mode().Perm() != 0700 {
		t.Errorf("Expected file permissions 0700, got %o", fileInfo.Mode().Perm())
	}

	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		t.Fatalf("Failed to stat directory: %v", err)
	}

	if dirInfo.Mode().Perm() != 0700 {
		t.Errorf("Expected directory permissions 0700, got %o", dirInfo.Mode().Perm())
	}
}

func TestCreator_EmptyPaths(t *testing.T) {
	creator := fs.NewCreator()

	err := creator.Create([]string{})
	if err != nil {
		t.Errorf("Create() with empty paths unexpected error = %v", err)
	}
}

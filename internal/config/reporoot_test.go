package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindProjectRoot_CIWorkspaceBoundaryHint(t *testing.T) {
	// Simulate CI.
	t.Setenv("CI", "true")

	// When running `go test`, the working directory is the package directory.
	// For internal/config, the repo root is ../..
	expectedRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("failed to resolve expected root: %v", err)
	}

	// Provide a CI boundary hint (must contain the start path).
	t.Setenv("FULMEN_WORKSPACE_ROOT", expectedRoot)

	root, err := FindProjectRoot()
	if err != nil {
		t.Fatalf("FindProjectRoot returned error: %v", err)
	}

	want, err := filepath.EvalSymlinks(expectedRoot)
	if err != nil {
		t.Fatalf("eval expected root: %v", err)
	}
	if filepath.Clean(root) != filepath.Clean(want) {
		t.Fatalf("expected root %q, got %q", want, root)
	}
}

func TestResolveCIBoundary_RejectsUnsafeHints(t *testing.T) {
	root := t.TempDir()
	start := filepath.Join(root, "repo")
	if err := os.Mkdir(start, 0o755); err != nil {
		t.Fatalf("mkdir start: %v", err)
	}
	outside := t.TempDir()
	file := filepath.Join(root, "not-a-dir")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if _, ok := resolveCIBoundary("relative/path", start); ok {
		t.Fatal("relative boundary was accepted")
	}
	if _, ok := resolveCIBoundary(outside, start); ok {
		t.Fatal("boundary that does not contain the start path was accepted")
	}
	if _, ok := resolveCIBoundary(file, start); ok {
		t.Fatal("file boundary was accepted")
	}

	got, ok := resolveCIBoundary(root, start)
	if !ok {
		t.Fatal("absolute directory containing the start path was rejected")
	}
	want, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatalf("eval root: %v", err)
	}
	if filepath.Clean(got) != filepath.Clean(want) {
		t.Fatalf("resolved boundary %q, want %q", got, want)
	}
}

func TestResolveCIBoundary_SymlinkResolvesToRealDirectory(t *testing.T) {
	real := t.TempDir()
	start := filepath.Join(real, "repo")
	if err := os.Mkdir(start, 0o755); err != nil {
		t.Fatalf("mkdir start: %v", err)
	}
	linkParent := t.TempDir()
	link := filepath.Join(linkParent, "workspace")
	if err := os.Symlink(real, link); err != nil {
		t.Fatalf("symlink: %v", err)
	}

	got, ok := resolveCIBoundary(link, start)
	if !ok {
		t.Fatal("symlink boundary containing the start path was rejected")
	}
	want, err := filepath.EvalSymlinks(real)
	if err != nil {
		t.Fatalf("eval real: %v", err)
	}
	if filepath.Clean(got) != filepath.Clean(want) {
		t.Fatalf("resolved boundary %q, want %q", got, want)
	}
}

func TestFindProjectRootFrom_MarkerlessBoundaryIsNotRoot(t *testing.T) {
	t.Setenv("CI", "true")
	for _, key := range []string{"GITHUB_WORKSPACE", "CI_PROJECT_DIR", "WORKSPACE", "PROJECT_ROOT", "GITHUB_ACTIONS"} {
		t.Setenv(key, "")
	}

	parent := t.TempDir()
	if err := os.WriteFile(filepath.Join(parent, "go.mod"), []byte("module example.com/marker\n"), 0o644); err != nil {
		t.Fatalf("write marker: %v", err)
	}
	markerless := filepath.Join(parent, "empty")
	start := filepath.Join(markerless, "sub")
	if err := os.MkdirAll(start, 0o755); err != nil {
		t.Fatalf("mkdir start: %v", err)
	}
	t.Setenv("FULMEN_WORKSPACE_ROOT", markerless)

	if _, err := FindProjectRootFrom(start); err == nil {
		t.Fatal("markerless bounded directory was accepted as a repository root")
	}
}

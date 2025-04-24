package converter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsImageFile(t *testing.T) {
	tests := []string{"foo.jpg", "bar.png", "baz.gif", "nope.txt"}
	want := []bool{true, true, true, false}
	for i, f := range tests {
		if isImageFile(f) != want[i] {
			t.Errorf("isImageFile(%s) = %v, want %v", f, isImageFile(f), want[i])
		}
	}
}

func TestConvertImage(t *testing.T) {
	settings := NewSettings()
	tmpDir := t.TempDir()
	src := filepath.Join("testdata", "sample.jpg")
	job := Job{SrcPath: src, RelPath: "sample.jpg"}
	err := ConvertImage(job, tmpDir, settings)
	if err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(tmpDir, "sample.webp")
	if _, err := os.Stat(out); err != nil {
		t.Errorf("output file not found: %v", err)
	}
}

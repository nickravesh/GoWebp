package converter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBatchConversion(t *testing.T) {
	settings := NewSettings()
	tmpDir := t.TempDir()
	jobs, err := FindImages([]string{"testdata"})
	if err != nil {
		t.Fatal(err)
	}
	logFile, _ := os.Create(filepath.Join(tmpDir, "conversion.log"))
	defer logFile.Close()
	count := 0
	WorkerPool(jobs, tmpDir, settings, logFile, func(done, total int) {
		count = done
	})
	if count != len(jobs) {
		t.Errorf("not all images converted: got %d, want %d", count, len(jobs))
	}
}

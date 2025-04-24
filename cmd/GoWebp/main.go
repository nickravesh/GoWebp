package main

import (
	"log"
	"os"
	"path/filepath"

	"GoWebp/internal/converter"
	"GoWebp/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	// "fyne.io/fyne/v2/storage"
)

func main() {
	// Detect working directory and create output dir
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory:", err)
	}
	outputDir := filepath.Join(cwd, "output-webp")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatal("Failed to create output directory:", err)
	}

	// Initialize settings (could be loaded from config)
	settings := converter.NewSettings()

	// Start Fyne app
	// a := app.New()
	a := app.NewWithID("com.nickravesh.GoWebp")
	w := a.NewWindow("GoWebp")

	// Set custom icon
	iconBytes, err := os.ReadFile("icon.png")
	if err == nil {
		resource := fyne.NewStaticResource("icon.png", iconBytes)
		w.SetIcon(resource)
	}

	ui.SetupUI(w, outputDir, settings)
	w.ShowAndRun()
}

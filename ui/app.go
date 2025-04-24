package ui

import (
	"GoWebp/internal/converter"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func SetupUI(w fyne.Window, outputDir string, settings *converter.Settings) {
	var selectedPaths []string
	progressBar := widget.NewProgressBar()
	logPane := widget.NewMultiLineEntry()
	logPane.SetMinRowsVisible(8)
	logPane.Disable()
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(500, 350))

	// Drag-and-drop area (OLD)
	// ddLabel := widget.NewLabel("Drag and drop images or folders here")
	// ddBox := container.NewVBox(ddLabel)
	// ddBox.Resize(fyne.NewSize(400, 100))

	// Drag-and-drop area (visual only)
	ddLabel := widget.NewLabel("Drag and drop images or folders anywhere on the window")

	// Enable file drop on the window
	w.SetOnDropped(func(pos fyne.Position, uris []fyne.URI) {
		for _, uri := range uris {
			selectedPaths = append(selectedPaths, uri.Path())
			logPane.SetText(logPane.Text + "\nAdded (dropped): " + uri.Path())
		}
	})

	// Browse button
	browseBtn := widget.NewButton("Add Image", func() {
		dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
			if uc != nil {
				selectedPaths = append(selectedPaths, uc.URI().Path())
				logPane.SetText(logPane.Text + "\nAdded: " + uc.URI().Path())
			}
		}, w)
	})

	// Select Directory button
	selectDirBtn := widget.NewButton("Add Directory", func() {
		dialog.ShowFolderOpen(func(lu fyne.ListableURI, err error) {
			if lu != nil {
				selectedPaths = append(selectedPaths, lu.Path())
				logPane.SetText(logPane.Text + "\nAdded directory: " + lu.Path())
			}
		}, w)
	})
	// Make both buttons expand equally and be responsive
	// browseBtnContainer := container.New(layout.NewStackLayout(), browseBtn)
	// selectDirBtnContainer := container.New(layout.NewStackLayout(), selectDirBtn)
	// browseRow := container.New(layout.NewGridLayoutWithColumns(2), browseBtnContainer, selectDirBtnContainer)

	// Settings button
	settingsBtn := widget.NewButton("Settings", func() {
		ShowSettingsDialog(w, settings)
	})

	// Start button
	startBtn := widget.NewButton("Start Conversion", func() {
		if len(selectedPaths) == 0 {
			dialog.ShowInformation("No files", "Please select files or folders.", w)
			return
		}
		logPane.SetText(logPane.Text + "\nProcessing files..., please wait.")
		jobs, err := converter.FindImages(selectedPaths)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		logFile, _ := os.Create(filepath.Join(outputDir, "conversion.log"))
		defer logFile.Close()
		progressBar.SetValue(0)
		total := len(jobs)
		go func() {
			done := 0
			converter.WorkerPool(jobs, outputDir, settings, logFile, func(d, t int) {
				done = d
				progressBar.SetValue(float64(done) / float64(total))
			})
			dialog.ShowInformation("Done", "Conversion complete!", w)
		}()
	})

	// Open output folder button
	openBtn := widget.NewButton("Open Output Folder", func() {
		go func() {
			path := outputDir
			if _, err := os.Stat(path); err == nil {
				if fyne.CurrentDevice().IsMobile() {
					dialog.ShowInformation("Not Supported", "Opening folders is not supported on mobile devices.", w)
				} else {
					var cmd string
					var args []string
					if os.Getenv("OS") == "Windows_NT" {
						cmd = "explorer"
						args = []string{path}
					} else {
						cmd = "xdg-open"
						args = []string{path}
					}
					if err := exec.Command(cmd, args...).Start(); err != nil {
						dialog.ShowError(err, w)
					}
				}
			} else {
				dialog.ShowError(err, w)
			}
		}()
	})

	// about button

	// about button
	aboutBtn := widget.NewButton("About", func() {
		dialog.ShowInformation(
			"About",
			"GoWebp\nImage to WebP Converter\nVersion 0.9.5 Beta Pre-release\n\nMade with ☕️ and ❤️ by ZED_ONE\n\nLicensed under GNU GPL v3\n\nGitHub: github.com/nickravesh/GoWebp",
			w,
		)
	})

	// aboutBtn := widget.NewButton("About", func() {
	// 	content := container.NewVBox(
	// 		widget.NewLabel("GoWebp\nImage to WebP Converter\nVersion 0.9.5 Beta Pre-release\n\nMade with ☕️ and ❤️ by ZED_ONE\n\nLicensed under GNU GPL v3"),
	// 		widget.NewHyperlink("GitHub: https://github.com/nickravesh/GoWebp",
	// 			parseURL("https://github.com/nickravesh/GoWebp")),
	// 	)
	// 	dialog.ShowCustom("About", "Close", content, w)
	// })

	w.SetContent(container.NewVBox(
		// ddBox,
		ddLabel,
		container.New(layout.NewGridLayoutWithColumns(3), browseBtn, selectDirBtn, settingsBtn),
		startBtn,
		progressBar,
		logPane,
		// container.NewBorder(nil, nil, nil, aboutBtn, openBtn),
		container.NewAdaptiveGrid(2, openBtn, aboutBtn),
	))
}

// Helper function to parse URL
func parseURL(urlStr string) *url.URL {
	u, _ := url.Parse(urlStr)
	return u
}

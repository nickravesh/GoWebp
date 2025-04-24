package ui

import (
	"fyne.io/fyne/v2/widget"
)

// NewProgressBar creates and returns a new progress bar widget.
func NewProgressBar() *widget.ProgressBar {
	bar := widget.NewProgressBar()
	bar.SetValue(0)
	return bar
}

package ui

import (
	"fyne.io/fyne/v2/widget"
)

// NewLogPane creates a multi-line, read-only log pane for status and logs.
func NewLogPane() *widget.Entry {
	logPane := widget.NewMultiLineEntry()
	logPane.SetMinRowsVisible(8)
	logPane.Disable()
	return logPane
}

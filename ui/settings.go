package ui

import (
	"GoWebp/internal/converter"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ShowSettingsDialog(w fyne.Window, settings *converter.Settings) {
	qualityEntry := widget.NewEntry()
	qualityEntry.SetText(strconv.Itoa(int(settings.Quality)))
	workersEntry := widget.NewEntry()
	workersEntry.SetText(strconv.Itoa(settings.MaxWorkers))
	overwriteCheck := widget.NewCheck("Allow overwrite", func(b bool) {
		settings.AllowOverwrite = b
	})
	overwriteCheck.SetChecked(settings.AllowOverwrite)

	formItems := []*widget.FormItem{
		widget.NewFormItem("Quality (0-100)", qualityEntry),
		widget.NewFormItem("Max Workers", workersEntry),
		widget.NewFormItem("", overwriteCheck),
	}
	dialog.ShowForm("Settings", "OK", "Cancel", formItems, func(b bool) {
		if b {
			if q, err := strconv.Atoi(qualityEntry.Text); err == nil {
				settings.Quality = float32(q)
			}
			if w, err := strconv.Atoi(workersEntry.Text); err == nil {
				settings.MaxWorkers = w
			}
		}
	}, w)
}

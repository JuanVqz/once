package ui

import (
	"context"

	tea "charm.land/bubbletea/v2"

	"github.com/basecamp/once/internal/docker"
)

const (
	backupsPathField = iota
	backupsAutoBackField
)

type SettingsFormBackups struct {
	settingsFormBase
	settings docker.ApplicationSettings
}

func NewSettingsFormBackups(app *docker.Application, lastResult *docker.OperationResult) *SettingsFormBackups {
	pathField := NewTextField("/path/to/backups")
	pathField.SetValue(app.Settings.Backup.Path)

	autoBackField := NewCheckboxField("Automatically create backups", app.Settings.Backup.AutoBack)

	m := &SettingsFormBackups{
		settingsFormBase: settingsFormBase{
			title: "Backups",
			form: NewForm("Done",
				FormItem{Label: "Backup location", Field: pathField},
				FormItem{Label: "Backups", Field: autoBackField},
			),
		},
		settings: app.Settings,
	}

	m.statusLine = func() string {
		return formatOperationStatus("backup", lastResult)
	}

	m.form.SetActionButton("Run backup now", func() tea.Msg {
		return settingsRunActionMsg{action: func() (string, error) {
			return "Backup complete", runBackup(app, pathField.Value())
		}}
	})
	m.form.OnSubmit(func() tea.Cmd {
		m.settings.Backup.Path = m.form.TextField(backupsPathField).Value()
		m.settings.Backup.AutoBack = m.form.CheckboxField(backupsAutoBackField).Checked()
		return func() tea.Msg { return SettingsSectionSubmitMsg{Settings: m.settings} }
	})
	m.form.OnCancel(func() tea.Cmd {
		return func() tea.Msg { return SettingsSectionCancelMsg{} }
	})

	return m
}

// Helpers

func runBackup(app *docker.Application, dir string) error {
	return app.BackupToFile(context.Background(), dir, app.BackupName())
}

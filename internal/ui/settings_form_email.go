package ui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/basecamp/once/internal/docker"
)

const (
	emailServerField = iota
	emailPortField
	emailUsernameField
	emailPasswordField
	emailFromField
)

type SettingsFormEmail struct {
	settingsFormBase
	settings docker.ApplicationSettings
}

func NewSettingsFormEmail(settings docker.ApplicationSettings) *SettingsFormEmail {
	serverField := NewTextField("smtp.example.com")
	serverField.SetValue(settings.SMTP.Server)

	portField := NewTextField("587")
	portField.SetCharLimit(5)
	portField.SetValue(settings.SMTP.Port)

	usernameField := NewTextField("user@example.com")
	usernameField.SetValue(settings.SMTP.Username)

	passwordField := NewTextField("password")
	passwordField.SetEchoPassword()
	passwordField.SetValue(settings.SMTP.Password)

	fromField := NewTextField("noreply@example.com")
	fromField.SetValue(settings.SMTP.From)

	m := &SettingsFormEmail{
		settingsFormBase: settingsFormBase{
			title: "Email",
			form: NewForm("Done",
				FormItem{Label: "SMTP Server", Field: serverField},
				FormItem{Label: "SMTP Port", Field: portField},
				FormItem{Label: "SMTP Username", Field: usernameField},
				FormItem{Label: "SMTP Password", Field: passwordField},
				FormItem{Label: "SMTP From", Field: fromField},
			),
		},
		settings: settings,
	}

	m.form.OnSubmit(func() tea.Cmd {
		m.settings.SMTP.Server = m.form.TextField(emailServerField).Value()
		m.settings.SMTP.Port = m.form.TextField(emailPortField).Value()
		m.settings.SMTP.Username = m.form.TextField(emailUsernameField).Value()
		m.settings.SMTP.Password = m.form.TextField(emailPasswordField).Value()
		m.settings.SMTP.From = m.form.TextField(emailFromField).Value()
		return func() tea.Msg { return SettingsSectionSubmitMsg{Settings: m.settings} }
	})
	m.form.OnCancel(func() tea.Cmd {
		return func() tea.Msg { return SettingsSectionCancelMsg{} }
	})

	return m
}

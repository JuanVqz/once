package ui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/basecamp/once/internal/docker"
)

const (
	appImageField = iota
	appHostnameField
	appTLSField
)

type SettingsFormApplication struct {
	settingsFormBase
	settings docker.ApplicationSettings
}

func NewSettingsFormApplication(settings docker.ApplicationSettings) *SettingsFormApplication {
	imageField := NewTextField("user/repo:tag")
	imageField.SetValue(settings.Image)

	hostnameField := NewTextField("app.example.com")
	hostnameField.SetValue(settings.Host)

	tlsField := NewCheckboxField("Enabled", !settings.DisableTLS)
	tlsField.SetDisabledWhen(func() (bool, string) {
		if docker.IsLocalhost(hostnameField.Value()) {
			return true, "Not available for localhost"
		}
		return false, ""
	})

	m := &SettingsFormApplication{
		settingsFormBase: settingsFormBase{
			title: "Application",
			form: NewForm("Done",
				FormItem{Label: "Image", Field: imageField, Required: true},
				FormItem{Label: "Hostname", Field: hostnameField, Required: true},
				FormItem{Label: "TLS", Field: tlsField},
			),
		},
		settings: settings,
	}

	m.form.OnSubmit(func() tea.Cmd {
		m.settings.Image = m.form.TextField(appImageField).Value()
		m.settings.Host = m.form.TextField(appHostnameField).Value()
		m.settings.DisableTLS = !m.form.CheckboxField(appTLSField).Checked()
		return func() tea.Msg { return SettingsSectionSubmitMsg{Settings: m.settings} }
	})
	m.form.OnCancel(func() tea.Cmd {
		return func() tea.Msg { return SettingsSectionCancelMsg{} }
	})

	return m
}

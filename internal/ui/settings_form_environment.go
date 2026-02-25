package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/basecamp/once/internal/docker"
)

type SettingsFormEnvironment struct {
	settingsFormBase
}

func NewSettingsFormEnvironment(settings docker.ApplicationSettings) *SettingsFormEnvironment {
	m := &SettingsFormEnvironment{
		settingsFormBase: settingsFormBase{
			title: "Environment",
			form:  NewForm("Done"),
		},
	}

	m.viewFn = func() string {
		placeholder := lipgloss.NewStyle().
			Foreground(Colors.Border).
			Italic(true).
			Render("(Environment variable editing coming soon)")

		return lipgloss.JoinVertical(lipgloss.Left,
			placeholder,
			"",
			m.form.View(),
		)
	}

	m.form.OnSubmit(func() tea.Cmd {
		return func() tea.Msg { return SettingsSectionCancelMsg{} }
	})
	m.form.OnCancel(func() tea.Cmd {
		return func() tea.Msg { return SettingsSectionCancelMsg{} }
	})

	return m
}

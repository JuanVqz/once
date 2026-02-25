package ui

import (
	"strconv"

	tea "charm.land/bubbletea/v2"

	"github.com/basecamp/once/internal/docker"
)

const (
	resourcesCPUField = iota
	resourcesMemoryField
)

type SettingsFormResources struct {
	settingsFormBase
	settings docker.ApplicationSettings
}

func NewSettingsFormResources(settings docker.ApplicationSettings) *SettingsFormResources {
	cpuField := NewTextField("e.g. 2")
	cpuField.SetCharLimit(10)
	cpuField.SetDigitsOnly(true)
	if settings.Resources.CPUs != 0 {
		cpuField.SetValue(strconv.Itoa(settings.Resources.CPUs))
	}

	memoryField := NewTextField("e.g. 512")
	memoryField.SetCharLimit(10)
	memoryField.SetDigitsOnly(true)
	if settings.Resources.MemoryMB != 0 {
		memoryField.SetValue(strconv.Itoa(settings.Resources.MemoryMB))
	}

	m := &SettingsFormResources{
		settingsFormBase: settingsFormBase{
			title: "Resources",
			form: NewForm("Done",
				FormItem{Label: "CPU Limit", Field: cpuField},
				FormItem{Label: "Memory Limit (MB)", Field: memoryField},
			),
		},
		settings: settings,
	}

	m.form.OnSubmit(func() tea.Cmd {
		m.settings.Resources.CPUs, _ = strconv.Atoi(m.form.TextField(resourcesCPUField).Value())
		m.settings.Resources.MemoryMB, _ = strconv.Atoi(m.form.TextField(resourcesMemoryField).Value())
		return func() tea.Msg { return SettingsSectionSubmitMsg{Settings: m.settings} }
	})
	m.form.OnCancel(func() tea.Cmd {
		return func() tea.Msg { return SettingsSectionCancelMsg{} }
	})

	return m
}

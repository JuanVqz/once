package ui

import tea "charm.land/bubbletea/v2"

type settingsFormBase struct {
	form       *Form
	title      string
	statusLine func() string
	viewFn     func() string
}

func (b *settingsFormBase) Title() string {
	return b.title
}

func (b *settingsFormBase) Init() tea.Cmd {
	return b.form.Init()
}

func (b *settingsFormBase) Update(msg tea.Msg) tea.Cmd {
	return b.form.Update(msg)
}

func (b *settingsFormBase) View() string {
	if b.viewFn != nil {
		return b.viewFn()
	}
	return b.form.View()
}

func (b *settingsFormBase) StatusLine() string {
	if b.statusLine != nil {
		return b.statusLine()
	}
	return ""
}

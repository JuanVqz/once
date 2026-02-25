package ui

import (
	"testing"

	"github.com/basecamp/gliff/tui"
	"github.com/stretchr/testify/assert"
)

func TestForm_FocusCycling(t *testing.T) {
	form := NewForm("Submit",
		FormItem{Label: "First", Field: NewTextField("first")},
		FormItem{Label: "Second", Field: NewTextField("second")},
	)
	assert.Equal(t, 0, form.Focused())

	formPressTab(form)
	assert.Equal(t, 1, form.Focused())

	formPressTab(form)
	assert.Equal(t, 2, form.Focused(), "submit button")

	formPressTab(form)
	assert.Equal(t, 3, form.Focused(), "cancel button")

	formPressTab(form)
	assert.Equal(t, 0, form.Focused(), "wraps to first field")
}

func TestForm_ShiftTabCycling(t *testing.T) {
	form := NewForm("Submit",
		FormItem{Label: "First", Field: NewTextField("first")},
		FormItem{Label: "Second", Field: NewTextField("second")},
	)

	formPressShiftTab(form)
	assert.Equal(t, 3, form.Focused(), "cancel button")

	formPressShiftTab(form)
	assert.Equal(t, 2, form.Focused(), "submit button")

	formPressShiftTab(form)
	assert.Equal(t, 1, form.Focused())

	formPressShiftTab(form)
	assert.Equal(t, 0, form.Focused())
}

func TestForm_EnterAdvancesFocus(t *testing.T) {
	form := NewForm("Submit",
		FormItem{Label: "First", Field: NewTextField("first")},
		FormItem{Label: "Second", Field: NewTextField("second")},
	)

	formPressEnter(form)
	assert.Equal(t, 1, form.Focused())

	formPressEnter(form)
	assert.Equal(t, 2, form.Focused(), "submit button")
}

func TestForm_SubmitAction(t *testing.T) {
	form := NewForm("Done",
		FormItem{Label: "Field", Field: NewTextField("val")},
	)
	submitted := false
	form.OnSubmit(func() tui.Cmd {
		submitted = true
		return nil
	})

	formPressTab(form)
	assert.Equal(t, 1, form.Focused(), "submit button")

	form.Update(keyMsg(tui.KeyEnter, 0))
	assert.True(t, submitted)
}

func TestForm_CancelAction(t *testing.T) {
	form := NewForm("Done",
		FormItem{Label: "Field", Field: NewTextField("val")},
	)
	cancelled := false
	form.OnCancel(func() tui.Cmd {
		cancelled = true
		return nil
	})

	formPressTab(form)
	formPressTab(form)
	assert.Equal(t, 2, form.Focused(), "cancel button")

	form.Update(keyMsg(tui.KeyEnter, 0))
	assert.True(t, cancelled)
}

func TestForm_NoFields(t *testing.T) {
	form := NewForm("Done")
	assert.Equal(t, 0, form.Focused(), "submit button")

	formPressTab(form)
	assert.Equal(t, 1, form.Focused(), "cancel button")

	formPressTab(form)
	assert.Equal(t, 0, form.Focused(), "wraps to submit")
}

func TestTextField_DigitsOnly(t *testing.T) {
	field := NewTextField("number")
	field.SetDigitsOnly(true)
	field.Focus()

	field.Update(runeMsg('5'))
	assert.Equal(t, "5", field.Value())

	field.Update(runeMsg('a'))
	assert.Equal(t, "5", field.Value(), "non-digit rejected")

	field.Update(runeMsg('3'))
	assert.Equal(t, "53", field.Value())
}

func TestCheckboxField_Toggle(t *testing.T) {
	field := NewCheckboxField("Enable", false)
	assert.False(t, field.Checked())

	field.Update(runeMsg(' '))
	assert.True(t, field.Checked())

	field.Update(runeMsg(' '))
	assert.False(t, field.Checked())
}

func TestCheckboxField_Render(t *testing.T) {
	field := NewCheckboxField("TLS", true)
	assert.Equal(t, "[✓] TLS", field.Render())

	field.Update(runeMsg(' '))
	assert.Equal(t, "[ ] TLS", field.Render())
}

func TestCheckboxField_DisabledWhen(t *testing.T) {
	disabled := true
	field := NewCheckboxField("TLS", false)
	field.SetDisabledWhen(func() (bool, string) {
		return disabled, "Not available"
	})

	field.Update(runeMsg(' '))
	assert.False(t, field.Checked(), "toggle ignored when disabled")
	assert.Equal(t, "Not available", field.Render())

	disabled = false
	field.Update(runeMsg(' '))
	assert.True(t, field.Checked(), "toggle works when enabled")
	assert.Equal(t, "[✓] TLS", field.Render())
}

func TestForm_FieldValuesAccessible(t *testing.T) {
	form := NewForm("Submit",
		FormItem{Label: "Name", Field: NewTextField("name")},
	)

	formTypeText(form, "hello")
	assert.Equal(t, "hello", form.TextField(0).Value())
}

func TestForm_ValidationBlocksSubmitWhenRequiredEmpty(t *testing.T) {
	form := NewForm("Submit",
		FormItem{Label: "Name", Field: NewTextField("name"), Required: true},
	)
	submitted := false
	form.OnSubmit(func() tui.Cmd {
		submitted = true
		return nil
	})

	formFocusSubmit(form)
	formPressEnter(form)

	assert.False(t, submitted)
	assert.True(t, form.HasError())
	assert.Equal(t, "Name is required", form.Error())
}

func TestForm_ValidationAllowsSubmitWhenRequiredFilled(t *testing.T) {
	form := NewForm("Submit",
		FormItem{Label: "Name", Field: NewTextField("name"), Required: true},
	)
	submitted := false
	form.OnSubmit(func() tui.Cmd {
		submitted = true
		return nil
	})

	formTypeText(form, "hello")
	formFocusSubmit(form)
	formPressEnter(form)

	assert.True(t, submitted)
	assert.False(t, form.HasError())
}

func TestForm_ValidationTreatsWhitespaceAsEmpty(t *testing.T) {
	form := NewForm("Submit",
		FormItem{Label: "Name", Field: NewTextField("name"), Required: true},
	)
	submitted := false
	form.OnSubmit(func() tui.Cmd {
		submitted = true
		return nil
	})

	formTypeText(form, "   ")
	formFocusSubmit(form)
	formPressEnter(form)

	assert.False(t, submitted)
	assert.Equal(t, "Name is required", form.Error())
}

func TestForm_ValidationErrorClearsOnInput(t *testing.T) {
	form := NewForm("Submit",
		FormItem{Label: "Name", Field: NewTextField("name"), Required: true},
	)
	form.OnSubmit(func() tui.Cmd { return nil })

	formFocusSubmit(form)
	formPressEnter(form)

	assert.True(t, form.HasError())
	assert.Equal(t, 0, form.Focused())

	formTypeText(form, "x")

	assert.False(t, form.HasError())
}

func TestForm_ValidationNonRequiredDoesNotBlock(t *testing.T) {
	form := NewForm("Submit",
		FormItem{Label: "Optional", Field: NewTextField("opt")},
		FormItem{Label: "Required", Field: NewTextField("req"), Required: true},
	)
	submitted := false
	form.OnSubmit(func() tui.Cmd {
		submitted = true
		return nil
	})

	formPressTab(form) // focus second field
	formTypeText(form, "filled")
	formFocusSubmit(form)
	formPressEnter(form)

	assert.True(t, submitted)
	assert.False(t, form.HasError())
}

func TestForm_ValidationOnClickSubmit(t *testing.T) {
	form := NewForm("Submit",
		FormItem{Label: "Name", Field: NewTextField("name"), Required: true},
	)
	submitted := false
	form.OnSubmit(func() tui.Cmd {
		submitted = true
		return nil
	})

	formClickSubmit(form)

	assert.False(t, submitted)
	assert.True(t, form.HasError())
}

func TestForm_ValidationFocusesFirstError(t *testing.T) {
	form := NewForm("Submit",
		FormItem{Label: "First", Field: NewTextField("first")},
		FormItem{Label: "Second", Field: NewTextField("second"), Required: true},
		FormItem{Label: "Third", Field: NewTextField("third"), Required: true},
	)
	form.OnSubmit(func() tui.Cmd { return nil })

	formFocusSubmit(form)
	formPressEnter(form)

	assert.Equal(t, 1, form.Focused(), "focused on first errored field")
	assert.Equal(t, "Second is required", form.Error())
}

// Helpers

func keyMsg(keyType tui.KeyType, r rune) tui.KeyMsg {
	return tui.KeyMsg{Key: tui.Key{Type: keyType, Rune: r}}
}

func runeMsg(r rune) tui.KeyMsg {
	return tui.KeyMsg{Key: tui.Key{Type: tui.KeyRune, Rune: r}}
}

func formPressTab(form *Form) {
	form.Update(keyMsg(tui.KeyTab, 0))
}

func formPressShiftTab(form *Form) {
	form.Update(keyMsg(tui.KeyShiftTab, 0))
}

func formPressEnter(form *Form) {
	form.Update(keyMsg(tui.KeyEnter, 0))
}

func formTypeText(form *Form, text string) {
	for _, r := range text {
		form.Update(runeMsg(r))
	}
}

func formFocusSubmit(form *Form) {
	for form.Focused() != form.submitIndex() {
		formPressTab(form)
	}
}

func formClickSubmit(form *Form) {
	form.Update(tui.MouseMsg{Type: tui.MousePress, Button: tui.MouseLeft, Target: "submit"})
}

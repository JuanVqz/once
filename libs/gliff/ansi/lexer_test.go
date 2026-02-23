package ansi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	t.Run("PlainText", func(t *testing.T) {
		l := NewLexer("Hello World")

		tok := l.Next()
		assert.Equal(t, TextToken, tok.Type)
		assert.Equal(t, "Hello World", tok.Text)

		tok = l.Next()
		assert.Equal(t, EOFToken, tok.Type)
	})

	t.Run("TextWithCSI", func(t *testing.T) {
		l := NewLexer("Hello\x1b[31mRed\x1b[0mWorld")

		tok := l.Next()
		assert.Equal(t, TextToken, tok.Type)
		assert.Equal(t, "Hello", tok.Text)

		tok = l.Next()
		assert.Equal(t, CSIToken, tok.Type)
		assert.Equal(t, "\x1b[31m", tok.Text)
		assert.Equal(t, "31", tok.Params)
		assert.Equal(t, byte('m'), tok.Final)

		tok = l.Next()
		assert.Equal(t, TextToken, tok.Type)
		assert.Equal(t, "Red", tok.Text)

		tok = l.Next()
		assert.Equal(t, CSIToken, tok.Type)
		assert.Equal(t, "\x1b[0m", tok.Text)
		assert.Equal(t, "0", tok.Params)
		assert.Equal(t, byte('m'), tok.Final)

		tok = l.Next()
		assert.Equal(t, TextToken, tok.Type)
		assert.Equal(t, "World", tok.Text)

		tok = l.Next()
		assert.Equal(t, EOFToken, tok.Type)
	})

	t.Run("MultipleCSISequences", func(t *testing.T) {
		l := NewLexer("\x1b[1;31m\x1b[44mText")

		tok := l.Next()
		assert.Equal(t, CSIToken, tok.Type)
		assert.Equal(t, "\x1b[1;31m", tok.Text)
		assert.Equal(t, "1;31", tok.Params)
		assert.Equal(t, byte('m'), tok.Final)

		tok = l.Next()
		assert.Equal(t, CSIToken, tok.Type)
		assert.Equal(t, "\x1b[44m", tok.Text)
		assert.Equal(t, "44", tok.Params)
		assert.Equal(t, byte('m'), tok.Final)

		tok = l.Next()
		assert.Equal(t, TextToken, tok.Type)
		assert.Equal(t, "Text", tok.Text)
	})

	t.Run("EmptyInput", func(t *testing.T) {
		l := NewLexer("")
		tok := l.Next()
		assert.Equal(t, EOFToken, tok.Type)
	})

	t.Run("OnlyCSI", func(t *testing.T) {
		l := NewLexer("\x1b[2J")

		tok := l.Next()
		assert.Equal(t, CSIToken, tok.Type)
		assert.Equal(t, "\x1b[2J", tok.Text)
		assert.Equal(t, "2", tok.Params)
		assert.Equal(t, byte('J'), tok.Final)

		tok = l.Next()
		assert.Equal(t, EOFToken, tok.Type)
	})

	t.Run("CSINoParams", func(t *testing.T) {
		l := NewLexer("\x1b[m")

		tok := l.Next()
		assert.Equal(t, CSIToken, tok.Type)
		assert.Equal(t, "\x1b[m", tok.Text)
		assert.Equal(t, "", tok.Params)
		assert.Equal(t, byte('m'), tok.Final)
	})

	t.Run("ESCOnly", func(t *testing.T) {
		l := NewLexer("\x1b")

		tok := l.Next()
		assert.Equal(t, TextToken, tok.Type)
		assert.Equal(t, "\x1b", tok.Text)

		tok = l.Next()
		assert.Equal(t, EOFToken, tok.Type)
	})

	t.Run("ESCWithNonCSI", func(t *testing.T) {
		l := NewLexer("\x1bM")

		tok := l.Next()
		assert.Equal(t, ESCToken, tok.Type)
		assert.Equal(t, "\x1bM", tok.Text)
		assert.Equal(t, byte('M'), tok.Final)

		tok = l.Next()
		assert.Equal(t, EOFToken, tok.Type)
	})
}

func TestParseCSIParams(t *testing.T) {
	t.Run("SingleParam", func(t *testing.T) {
		params := ParseCSIParams("31", nil)
		assert.Equal(t, []int{31}, params)
	})

	t.Run("MultipleParams", func(t *testing.T) {
		params := ParseCSIParams("1;31;44", nil)
		assert.Equal(t, []int{1, 31, 44}, params)
	})

	t.Run("Empty", func(t *testing.T) {
		params := ParseCSIParams("", nil)
		assert.Equal(t, []int{0}, params)
	})

	t.Run("RGBColor", func(t *testing.T) {
		params := ParseCSIParams("38;2;255;128;64", nil)
		assert.Equal(t, []int{38, 2, 255, 128, 64}, params)
	})

	t.Run("256Color", func(t *testing.T) {
		params := ParseCSIParams("38;5;196", nil)
		assert.Equal(t, []int{38, 5, 196}, params)
	})

	t.Run("WithEmptyParts", func(t *testing.T) {
		params := ParseCSIParams("1;;31", nil)
		assert.Equal(t, []int{1, 0, 31}, params)
	})

	t.Run("WithBuffer", func(t *testing.T) {
		var buf [16]int
		params := ParseCSIParams("1;31;44", buf[:])
		assert.Equal(t, []int{1, 31, 44}, params)
	})
}

func TestParseCSIParam(t *testing.T) {
	assert.Equal(t, 0, ParseCSIParam(""))
	assert.Equal(t, 42, ParseCSIParam("42"))
	assert.Equal(t, 123, ParseCSIParam("123"))
}

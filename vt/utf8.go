package vt

import (
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/wcwidth"
)

// handleUtf8 handles a UTF-8 characters.
func (t *Terminal) handleUtf8(seq ansi.Sequence) {
	var width int
	var content string
	switch seq := seq.(type) {
	case ansi.Rune:
		width = wcwidth.RuneWidth(rune(seq))
		content = string(seq)
	case ansi.Grapheme:
		width = seq.Width
		content = seq.Cluster
	}

	var autowrap bool
	x, y := t.scr.CursorPosition()
	if t.isModeSet(ansi.AutoWrapMode) {
		autowrap = true
	}

	// Handle wide chars at the edge - wrap them entirely
	if autowrap && x+width > t.scr.Width() {
		x = 0
		y++
		// Only scroll if we're past the last line
		if y >= t.scr.Height() {
			t.scr.ScrollUp(1)
			y = t.scr.Height() - 1
		}
	}

	// Handle character set mappings
	if len(content) == 1 {
		var charset CharSet
		c := content[0]
		if t.gsingle > 0 && t.gsingle < 4 {
			charset = t.charsets[t.gsingle]
			t.gsingle = 0
		} else if c < 128 {
			charset = t.charsets[t.gl]
		} else {
			charset = t.charsets[t.gr]
		}

		if charset != nil {
			if r, ok := charset[c]; ok {
				content = r
			}
		}
	}

	cell := &Cell{
		Style:   t.scr.cur.Pen,
		Link:    Link{}, // TODO: Link support
		Content: content,
		Width:   width,
	}

	t.scr.SetCell(x, y, cell)

	// TODO: Is this correct?
	t.scr.setCursor(x+width, y, true)
}

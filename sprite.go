// Copyright 2016 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import "strings"

type Sprite struct {
	Block
	Text        string
	TextFgColor Attribute
	TextBgColor Attribute
	WrapLength  int // words wrap limit. Note it may not work properly with multi-width char
	Alpha       rune
}

// NewPar returns a new *Par with given text as its content.
func NewSprite(s string) *Sprite {
	block := NewBlock()
	block.Width, block.Height = getWidthHeight(s)
	block.Border = false
	block.BorderLabel = ""

	return &Sprite{
		Block:       *block,
		Text:        s,
		TextFgColor: ThemeAttr("par.text.fg"),
		TextBgColor: ThemeAttr("par.text.bg"),
		WrapLength:  0,
		Alpha:       'x',
	}
}

func getWidthHeight(s string) (int, int) {
	var width int
	var height int

	for _, line := range strings.Split(s, "\n") {
		if len(line) > width {
			width = len(line)
		}
		height += 1
	}
	return width, height
}

// Buffer implements Bufferer interface.
func (s *Sprite) Buffer() Buffer {
	buf := s.Block.Buffer()

	fg, bg := s.TextFgColor, s.TextBgColor
	cs := DefaultTxBuilder.Build(s.Text, fg, bg)

	y, x, n := 0, 0, 0
	for y < s.innerArea.Dy() && n < len(cs) {
		w := cs[n].Width()

		if cs[n].Ch == s.Alpha {
			cs[n].Visible = false
		}

		if cs[n].Ch == '\n' || x+w > s.innerArea.Dx() {
			y++
			x = 0 // set x = 0
			if cs[n].Ch == '\n' {
				n++
			}
			continue
		}

		buf.Set(s.innerArea.Min.X+x, s.innerArea.Min.Y+y, cs[n])

		n++
		x += w
	}

	return buf
}

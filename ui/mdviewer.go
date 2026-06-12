package ui

import (
	"strings"

	"git.sr.ht/~rockorager/vaxis"
)

type markdownViewer struct {
	active bool
	title  string
	url    string
	lines  []StyledString
	scroll int
}

func (mv *markdownViewer) open(title, url, content string, width int) {
	mv.title = title
	mv.url = url
	mv.lines = renderMarkdown(content, width)
	mv.scroll = 0
	mv.active = true
}

func (mv *markdownViewer) close() {
	mv.active = false
	mv.lines = nil
}

func (mv *markdownViewer) scrollUp(n int) {
	mv.scroll -= n
	if mv.scroll < 0 {
		mv.scroll = 0
	}
}

func (mv *markdownViewer) scrollDown(n, innerH int) {
	mv.scroll += n
	max := len(mv.lines) - innerH
	if max < 0 {
		max = 0
	}
	if mv.scroll > max {
		mv.scroll = max
	}
}

func renderMarkdown(content string, width int) []StyledString {
	var out []StyledString
	inCode := false
	rawLines := strings.Split(content, "\n")

	add := func(s StyledString) {
		out = append(out, wrapStyledLine(s, width)...)
	}

	for _, raw := range rawLines {
		if strings.HasPrefix(raw, "```") {
			inCode = !inCode
			continue
		}
		if inCode {
			add(Styled("  "+raw, vaxis.Style{Foreground: vaxis.IndexColor(10)}))
			continue
		}
		switch {
		case strings.HasPrefix(raw, "# "):
			add(Styled(raw[2:], vaxis.Style{Foreground: vaxis.IndexColor(14), Attribute: vaxis.AttrBold}))
		case strings.HasPrefix(raw, "## "):
			add(Styled(raw[3:], vaxis.Style{Foreground: vaxis.IndexColor(11), Attribute: vaxis.AttrBold}))
		case strings.HasPrefix(raw, "### "):
			add(Styled(raw[4:], vaxis.Style{Attribute: vaxis.AttrBold}))
		case strings.HasPrefix(raw, "#### ") || strings.HasPrefix(raw, "##### ") || strings.HasPrefix(raw, "###### "):
			idx := strings.Index(raw, " ")
			add(Styled(raw[idx+1:], vaxis.Style{Attribute: vaxis.AttrBold}))
		case strings.TrimRight(raw, "-") == "" && len(strings.TrimSpace(raw)) >= 3 && strings.TrimSpace(raw)[0] == '-':
			sep := width
			if sep < 1 {
				sep = 1
			}
			out = append(out, PlainString(strings.Repeat("─", sep)))
		case strings.HasPrefix(raw, "- ") || strings.HasPrefix(raw, "* "):
			add(parseMdInline("• "+raw[2:]))
		case strings.HasPrefix(raw, "  - ") || strings.HasPrefix(raw, "  * "):
			add(parseMdInline("  • "+raw[4:]))
		default:
			add(parseMdInline(raw))
		}
	}
	return out
}

// sliceStyledBytes returns the sub-StyledString from startByte to endByte,
// carrying over the style in effect at startByte.
func sliceStyledBytes(s StyledString, startByte, endByte int) StyledString {
	sub := s.string[startByte:endByte]
	if len(sub) == 0 {
		return PlainString("")
	}
	var inherit vaxis.Style
	var newStyles []rangedStyle
	for _, rs := range s.styles {
		if rs.Start < startByte {
			inherit = rs.Style
		} else if rs.Start < endByte {
			newStyles = append(newStyles, rangedStyle{Start: rs.Start - startByte, Style: rs.Style})
		}
	}
	if len(newStyles) == 0 || newStyles[0].Start != 0 {
		newStyles = append([]rangedStyle{{Start: 0, Style: inherit}}, newStyles...)
	}
	return StyledString{string: sub, styles: newStyles}
}

// wrapStyledLine splits s into lines of at most maxWidth runes, breaking at spaces.
func wrapStyledLine(s StyledString, maxWidth int) []StyledString {
	if maxWidth <= 0 {
		return []StyledString{s}
	}
	sr := []rune(s.string)
	if len(sr) <= maxWidth {
		return []StyledString{s}
	}
	var result []StyledString
	runeStart := 0
	for runeStart < len(sr) {
		remaining := len(sr) - runeStart
		if remaining <= maxWidth {
			byteStart := len(string(sr[:runeStart]))
			result = append(result, sliceStyledBytes(s, byteStart, len(s.string)))
			break
		}
		// Find last space within maxWidth of runeStart.
		breakAt := runeStart + maxWidth
		spaceAt := -1
		for i := breakAt; i > runeStart; i-- {
			if sr[i] == ' ' {
				spaceAt = i
				break
			}
		}
		lineEnd := breakAt
		nextStart := breakAt
		if spaceAt > runeStart {
			lineEnd = spaceAt
			nextStart = spaceAt + 1
		}
		byteStart := len(string(sr[:runeStart]))
		byteEnd := len(string(sr[:lineEnd]))
		result = append(result, sliceStyledBytes(s, byteStart, byteEnd))
		runeStart = nextStart
	}
	return result
}

func parseMdInline(s string) StyledString {
	var sb StyledStringBuilder
	base := vaxis.Style{}
	sb.SetStyle(base)
	bold := false
	code := false
	i := 0
	for i < len(s) {
		switch {
		case !code && i+1 < len(s) && s[i] == '*' && s[i+1] == '*':
			bold = !bold
			if bold {
				sb.SetStyle(vaxis.Style{Attribute: vaxis.AttrBold})
			} else {
				sb.SetStyle(base)
			}
			i += 2
		case s[i] == '`':
			code = !code
			if code {
				sb.SetStyle(vaxis.Style{Foreground: vaxis.IndexColor(10)})
			} else {
				sb.SetStyle(base)
			}
			i++
		default:
			j := i + 1
			for j < len(s) {
				if s[j] == '`' {
					break
				}
				if s[j] == '*' && j+1 < len(s) && s[j+1] == '*' {
					break
				}
				j++
			}
			sb.WriteString(s[i:j])
			i = j
		}
	}
	return sb.StyledString()
}

func (ui *UI) drawMarkdownViewer(vx *Vaxis) {
	w, h := vx.window.Size()

	x0, y0 := 2, 1
	bw, bh := w-4, h-2
	if bw < 10 || bh < 4 {
		return
	}

	// Use explicit dark background so the overlay is always visible.
	bgColor := vaxis.IndexColor(235)
	fgColor := vaxis.IndexColor(255)
	bgSt := vaxis.Style{Background: bgColor, Foreground: fgColor}
	borderSt := vaxis.Style{Background: bgColor, Foreground: vaxis.IndexColor(244)}

	// Fill entire box with background.
	for y := y0; y < y0+bh; y++ {
		for x := x0; x < x0+bw; x++ {
			setCell(vx, x, y, ' ', bgSt)
		}
	}

	// top border
	setCell(vx, x0, y0, '┌', borderSt)
	setCell(vx, x0+bw-1, y0, '┐', borderSt)
	for x := x0 + 1; x < x0+bw-1; x++ {
		setCell(vx, x, y0, '─', borderSt)
	}
	// Title on the left of the top border.
	titleText := " " + ui.mdv.title + " "
	tx := x0 + 2
	printString(vx, &tx, y0, Styled(titleText, vaxis.Style{Background: bgColor, Foreground: vaxis.IndexColor(15), Attribute: vaxis.AttrBold}))

	// Commands hint on the right of the top border.
	cmds := " ↑/↓ · PgUp/PgDn · Esc close "
	if ui.mdv.url != "" {
		cmds = " ↑/↓ · PgUp/PgDn · b browser · Esc close "
	}
	cmdsRunes := []rune(cmds)
	cx := x0 + bw - 1 - len(cmdsRunes)
	if cx > tx+1 {
		printString(vx, &cx, y0, Styled(cmds, borderSt))
	}

	// bottom border (no hint needed — already in header)
	setCell(vx, x0, y0+bh-1, '└', borderSt)
	setCell(vx, x0+bw-1, y0+bh-1, '┘', borderSt)
	for x := x0 + 1; x < x0+bw-1; x++ {
		setCell(vx, x, y0+bh-1, '─', borderSt)
	}

	// side borders
	for y := y0 + 1; y < y0+bh-1; y++ {
		setCell(vx, x0, y, '│', borderSt)
		setCell(vx, x0+bw-1, y, '│', borderSt)
	}

	// scrollbar
	innerH := bh - 2
	total := len(ui.mdv.lines)
	if total > innerH {
		barH := innerH * innerH / total
		if barH < 1 {
			barH = 1
		}
		barY := y0 + 1 + (ui.mdv.scroll*(innerH-barH))/(total-innerH)
		for y := y0 + 1; y < y0+bh-1; y++ {
			ch := '░'
			if y >= barY && y < barY+barH {
				ch = '█'
			}
			setCell(vx, x0+bw-1, y, ch, borderSt)
		}
	}

	// content lines
	contentX := x0 + 1
	contentY := y0 + 1
	innerW := bw - 2 // subtract left border; right border/scrollbar already excluded
	xMax := contentX + innerW - 1
	for i := 0; i < innerH; i++ {
		li := i + ui.mdv.scroll
		if li >= len(ui.mdv.lines) {
			break
		}
		x := contentX
		line := injectBackground(ui.mdv.lines[li], bgColor)
		printStringClamped(vx, &x, contentY+i, xMax, line)
	}
}

// printStringClamped is like printString but stops rendering at xMax.
func printStringClamped(vx *Vaxis, x *int, y, xMax int, s StyledString) {
	var st vaxis.Style
	nextStyles := s.styles
	i := 0
	sr := []rune(s.string)
	for len(sr) > 0 {
		if *x >= xMax {
			break
		}
		if len(nextStyles) > 0 && nextStyles[0].Start == i {
			st = nextStyles[0].Style
			nextStyles = nextStyles[1:]
		}
		dx, di := printCluster(vx, *x, y, xMax, sr, st)
		if di == 0 {
			break
		}
		*x += dx
		i += len(string(sr[:di]))
		sr = sr[di:]
	}
}

// injectBackground returns a copy of s with bg applied to all style ranges
// that have no explicit background set.
func injectBackground(s StyledString, bg vaxis.Color) StyledString {
	if len(s.styles) == 0 {
		return Styled(s.string, vaxis.Style{Background: bg})
	}
	styles := make([]rangedStyle, len(s.styles))
	copy(styles, s.styles)
	for i := range styles {
		if styles[i].Style.Background == ColorDefault {
			styles[i].Style.Background = bg
		}
	}
	// Ensure position 0 has a style with background.
	if styles[0].Start != 0 {
		styles = append([]rangedStyle{{Start: 0, Style: vaxis.Style{Background: bg}}}, styles...)
	}
	return StyledString{string: s.string, styles: styles}
}

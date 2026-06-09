package main

import (
	"fmt"
	"image/color"
	"log"
	"strings"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		log.Fatal(err)
	}
}

type errMsg error

type model struct {
	textArea textarea.Model
	err      error
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Once upon a time..."
	ta.SetVirtualCursor(false)
	ta.SetStyles(textarea.DefaultStyles(true))
	ta.Focus()

	return model{
		textArea: ta,
	}
}

// Init implements [tea.Model].
func (m model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		tea.RequestBackgroundColor,
	)
}

// Update implements [tea.Model].
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		// Update the styling now that we know the background color.
		m.textArea.SetStyles(textarea.DefaultStyles(msg.IsDark()))

	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			if m.textArea.Focused() {
				m.textArea.Blur()
			}
		case "ctrl+c":
			return m, tea.Quit
		default:
			if !m.textArea.Focused() {
				cmd = m.textArea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textArea, cmd = m.textArea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (model) headerView() string {
	return "This text area is styled.\n"
}

// View implements [tea.Model].
func (m model) View() tea.View {
	const (
		footer = "\n(ctrl+c to quit)\n"
	)

	var c *tea.Cursor
	if !m.textArea.VirtualCursor() {
		c = m.textArea.Cursor()

		if c != nil {
			offset := lipgloss.Height(m.headerView())
			c.Y += offset
		}
	}

	f := strings.Join([]string{
		m.headerView(),
		m.textArea.View(),
		footer,
	}, "\n")

	v := tea.NewView(f)
	v.Cursor = c
	return v
}

type css map[string]any

func styleToCSS(style lipgloss.Style) css {
	var (
		align           = style.GetAlign()
		alignHorizontal = style.GetAlignHorizontal()
		alignVertical   = style.GetAlignVertical()
		background      = style.GetBackground()
		blink           = style.GetBlink()
		bold            = style.GetBold()

		// border, borderTop, borderRight, borderBottom, borderLeft = style.GetBorder()
		borderForegroundBlend       = style.GetBorderForegroundBlend()
		borderForegroundBlendOffset = style.GetBorderForegroundBlendOffset()
		borderStyle                 = style.GetBorderStyle()
		borderBottom                = style.GetBorderBottom()
		borderBottomBackground      = style.GetBorderBottomBackground()
		borderBottomForeground      = style.GetBorderBottomForeground()
		borderBottomSize            = style.GetBorderBottomSize()
		borderLeft                  = style.GetBorderLeft()
		borderLeftBackground        = style.GetBorderLeftBackground()
		borderLeftForeground        = style.GetBorderLeftForeground()
		borderLeftSize              = style.GetBorderLeftSize()
		borderRight                 = style.GetBorderRight()
		borderRightBackground       = style.GetBorderRightBackground()
		borderRightForeground       = style.GetBorderRightForeground()
		borderRightSize             = style.GetBorderRightSize()
		borderTop                   = style.GetBorderTop()
		borderTopBackground         = style.GetBorderTopBackground()
		borderTopForeground         = style.GetBorderTopForeground()
		borderTopSize               = style.GetBorderTopSize()

		colorWhitespace        = style.GetColorWhitespace()
		faint                  = style.GetFaint()
		foreground             = style.GetForeground()
		frameSizeX, frameSizeY = style.GetFrameSize()
		height                 = style.GetHeight()
		horizontalBorderSize   = style.GetHorizontalBorderSize()
		horizontalFrameSize    = style.GetHorizontalFrameSize()
		horizontalMargins      = style.GetHorizontalMargins()
		horizontalPadding      = style.GetHorizontalPadding()
		hyperlink, params      = style.GetHyperlink()
		inline                 = style.GetInline()
		italic                 = style.GetItalic()

		marginTop, marginRight, marginBottom, marginLeft = style.GetMargin()
		marginChar                                       = style.GetMarginChar()
		// _ = style.GetMarginBottom()
		// _ = style.GetMarginLeft()
		// _ = style.GetMarginRight()
		// _ = style.GetMarginTop()

		maxHeight = style.GetMaxHeight()
		maxWidth  = style.GetMaxWidth()

		paddingTop, paddingRight, paddingBottom, paddingLeft = style.GetPadding()
		paddingChar                                          = style.GetPaddingChar()
		// _ = style.GetPaddingBottom()
		// _ = style.GetPaddingLeft()
		// _ = style.GetPaddingRight()
		// _ = style.GetPaddingTop()

		reverse             = style.GetReverse()
		strikethrough       = style.GetStrikethrough()
		strikethroughSpaces = style.GetStrikethroughSpaces()
		tabWidth            = style.GetTabWidth()
		// transform = style.GetTransform()
		underline          = style.GetUnderline()
		underlineColor     = style.GetUnderlineColor()
		underlineSpaces    = style.GetUnderlineSpaces()
		underlineStyle     = style.GetUnderlineStyle()
		verticalBorderSize = style.GetVerticalBorderSize()
		verticalFrameSize  = style.GetVerticalFrameSize()
		verticalMargins    = style.GetVerticalMargins()
		verticalPadding    = style.GetVerticalPadding()
		width              = style.GetWidth()
	)

	return css{
		"align":                          align,
		"align-horizontal":               alignHorizontal,
		"align-vertical":                 alignVertical,
		"background":                     asHex(background),
		"blink":                          blink,
		"bold":                           bold,
		"border-foreground-blend":        asHexes(borderForegroundBlend),
		"border-foreground-blend-offset": borderForegroundBlendOffset,
		"border-style":                   borderStyle,
		"border-bottom":                  borderBottom,
		"border-bottom-background":       borderBottomBackground,
		"border-bottom-foreground":       borderBottomForeground,
		"border-bottom-size":             borderBottomSize,
		"border-left":                    borderLeft,
		"border-left-background":         borderLeftBackground,
		"border-left-foreground":         borderLeftForeground,
		"border-left-size":               borderLeftSize,
		"border-right":                   borderRight,
		"border-right-background":        borderRightBackground,
		"border-right-foreground":        borderRightForeground,
		"border-right-size":              borderRightSize,
		"border-top":                     borderTop,
		"border-top-background":          borderTopBackground,
		"border-top-foreground":          borderTopForeground,
		"border-top-size":                borderTopSize,
		"color-whitespace":               colorWhitespace,
	}
}

func asHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02X%02X%02X", r>>8, g>>8, b>>8)
}

func asHexes(cs []color.Color) []string {
	var colors []string
	for _, c := range cs {
		colors = append(colors, asHex(c))
	}
	return colors
}

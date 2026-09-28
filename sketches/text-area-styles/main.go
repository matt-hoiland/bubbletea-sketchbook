package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"
	"time"

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
	quitting bool
	err      error
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Once upon a time..."
	ta.SetVirtualCursor(false)
	styles := textarea.DefaultDarkStyles()
	styles.Focused.Base.
		Foreground(lipgloss.Color("#f48c24")).
		Background(lipgloss.Color("#50246c")).
		Bold(true)
	ta.SetStyles(styles)
	ta.Focus()

	return model{
		textArea: ta,
	}
}

// Init implements [tea.Model].
func (m model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
	)
}

// Update implements [tea.Model].
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			if m.textArea.Focused() {
				m.textArea.Blur()
			}
		case "ctrl+c":
			m.quitting = true
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

	if m.quitting {
		styles := m.textArea.Styles()
		data, err := json.MarshalIndent(stylesToCSS(styles).dropZeroes(), "", "  ")
		if err != nil {
			panic(err)
		}
		if err := os.WriteFile("styles.json", data, 0644); err != nil {
			panic(err)
		}
	}

	return v
}

type css map[string]any

func (m css) dropZeroes() css {
	filtered := make(css)

	for key, value := range m {
		switch value := value.(type) {
		case css:
			filteredValue := value.dropZeroes()
			if len(filteredValue) > 0 {
				filtered[key] = filteredValue
			}
		case bool:
			if value {
				filtered[key] = value
			}
		case int:
			if value != 0 {
				filtered[key] = value
			}
		case int32: //rune
			if value != 0 {
				filtered[key] = fmt.Sprintf("%c", value)
			}
		case uint8:
			if value != 0 {
				filtered[key] = value
			}
		case string:
			if len(value) > 0 && value != "#000000" {
				filtered[key] = value
			}
		case []string:
			if len(value) > 0 {
				filtered[key] = value
			}
		case lipgloss.Border:
			var borderless lipgloss.Border
			if value != borderless {
				filtered[key] = value
			}
		case lipgloss.Position:
			if value != 0 {
				filtered[key] = value
			}
		case lipgloss.NoColor:
			continue
		case tea.CursorShape:
			if value != 0 {
				filtered[key] = value
			}
		case time.Duration:
			if value != 0 {
				filtered[key] = value
			}
		default:
			panic(fmt.Errorf("unknown type: key \"%s\", value (%T) %v\n", key, value, value))
		}
	}
	return filtered
}

func stylesToCSS(styles textarea.Styles) css {
	return css{
		"focused": styleStateToCSS(styles.Focused),
		"blurred": styleStateToCSS(styles.Blurred),
		"cursor": css{
			"color":       asHex(styles.Cursor.Color),
			"shape":       styles.Cursor.Shape,
			"blink":       styles.Cursor.Blink,
			"blink-speed": styles.Cursor.BlinkSpeed,
		},
	}
}

func styleStateToCSS(styleState textarea.StyleState) css {
	return css{
		"base":               styleToCSS(styleState.Base),
		"text":               styleToCSS(styleState.Text),
		"line-number":        styleToCSS(styleState.LineNumber),
		"cursor-line-number": styleToCSS(styleState.CursorLineNumber),
		"cursor-line":        styleToCSS(styleState.CursorLine),
		"end-of-buffer":      styleToCSS(styleState.EndOfBuffer),
		"placeholder":        styleToCSS(styleState.Placeholder),
		"prompt":             styleToCSS(styleState.Prompt),
	}
}

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

		colorWhitespace            = style.GetColorWhitespace()
		faint                      = style.GetFaint()
		foreground                 = style.GetForeground()
		frameSizeX, frameSizeY     = style.GetFrameSize()
		height                     = style.GetHeight()
		horizontalBorderSize       = style.GetHorizontalBorderSize()
		horizontalFrameSize        = style.GetHorizontalFrameSize()
		horizontalMargins          = style.GetHorizontalMargins()
		horizontalPadding          = style.GetHorizontalPadding()
		hyperlink, hyperlinkParams = style.GetHyperlink()
		inline                     = style.GetInline()
		italic                     = style.GetItalic()

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
		"faint":                          faint,
		"foreground":                     asHex(foreground),
		"frame-size-x":                   frameSizeX,
		"frame-size-y":                   frameSizeY,
		"height":                         height,
		"horizontal-border-size":         horizontalBorderSize,
		"horizontal-frame-size":          horizontalFrameSize,
		"horizontal-margins":             horizontalMargins,
		"horizontal-padding":             horizontalPadding,
		"hyperlink":                      hyperlink,
		"hyperlink-params":               hyperlinkParams,
		"inline":                         inline,
		"italic":                         italic,
		"margin-top":                     marginTop,
		"margin-right":                   marginRight,
		"margin-bottom":                  marginBottom,
		"margin-left":                    marginLeft,
		"margin-char":                    marginChar,
		"max-height":                     maxHeight,
		"max-width":                      maxWidth,
		"padding-top":                    paddingTop,
		"padding-right":                  paddingRight,
		"padding-bottom":                 paddingBottom,
		"padding-left":                   paddingLeft,
		"padding-char":                   paddingChar,
		"reverse":                        reverse,
		"strikethrough":                  strikethrough,
		"strikethrough-spaces":           strikethroughSpaces,
		"tab-width":                      tabWidth,
		"underline":                      underline,
		"underline-color":                underlineColor,
		"underline-spaces":               underlineSpaces,
		"underline-style":                underlineStyle,
		"vertical-border-size":           verticalBorderSize,
		"vertical-frame-size":            verticalFrameSize,
		"vertical-margins":               verticalMargins,
		"vertical-padding":               verticalPadding,
		"width":                          width,
	}
}

func asHex(c color.Color) string {
	if _, ok := c.(lipgloss.NoColor); ok {
		return ""
	}
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

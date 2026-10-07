package window

import (
	"fmt"
	"image/color"
	"io"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/paginator"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	catppuccingo "github.com/catppuccin/go"
)

type Model struct {
	keys keyMap
	help help.Model

	list list.Model

	width, height int

	title string
}

type Option func(*conf)

type conf struct {
	title string
}

func New(opts ...Option) *Model {
	c := conf{
		title: "program",
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&c)
		}
	}

	ll := list.New([]list.Item{
		Task{Incomplete, "call Mom and Dad"},
		Task{InProgress, "go to bed by 10 PM"},
		Task{Complete, "watch PBS NewsHour"},
		Task{Cancelled, "workout"},
		Task{Incomplete, "call Mom and Dad"},
		Task{Incomplete, "go to bed by 10 PM"},
		Task{Complete, "watch PBS NewsHour"},
		Task{Cancelled, "workout"},
		Task{Incomplete, "call Mom and Dad"},
		Task{Incomplete, "go to bed by 10 PM"},
		Task{Complete, "watch PBS NewsHour"},
		Task{Cancelled, "workout"},
		Task{Incomplete, "call Mom and Dad"},
		Task{Incomplete, "go to bed by 10 PM"},
		Task{Complete, "watch PBS NewsHour"},
		Task{Cancelled, "workout"},
		Task{Incomplete, "call Mom and Dad"},
		Task{Incomplete, "go to bed by 10 PM"},
		Task{Complete, "watch PBS NewsHour"},
		Task{Cancelled, "workout"},
		Task{Incomplete, "call Mom and Dad"},
		Task{Incomplete, "go to bed by 10 PM"},
		Task{Complete, "watch PBS NewsHour"},
		Task{Cancelled, "workout"},
	}, TaskDelegate{}, 0, 0)

	// ll.SetShowPagination(false)
	ll.InfiniteScrolling = true
	ll.Paginator.Type = paginator.Arabic

	m := &Model{
		keys: keys,
		help: help.New(),

		list: ll,

		title: c.title,
	}

	return m
}

func WithTitle(title string) Option {
	return func(c *conf) {
		c.title = title
	}
}

type status int

const (
	Incomplete status = iota
	InProgress
	Complete
	Cancelled
)

type Task struct {
	status     status
	imperative string
}

func (t Task) Title() string { return t.imperative }
func (t Task) Description() string {
	switch t.status {
	case Incomplete:
		return "incomplete"
	case InProgress:
		return "in progress"
	case Complete:
		return "complete"
	case Cancelled:
		return "cancelled"
	default:
		return "default"
	}
}
func (t Task) FilterValue() string { return t.imperative }

type TaskDelegate struct{}

func (TaskDelegate) Height() int                         { return 1 }
func (TaskDelegate) Spacing() int                        { return 0 }
func (TaskDelegate) Update(tea.Msg, *list.Model) tea.Cmd { return nil }

func (TaskDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	task, ok := item.(Task)
	if !ok {
		task = Task{status: Incomplete, imperative: "unknown task"}
	}

	background := catppuccingo.Mocha.Surface0()

	var (
		highlight = false
		box       = lipgloss.NewStyle().PaddingLeft(2)
	)
	if index == m.Index() {
		highlight = true
		box = lipgloss.NewStyle().
			Width(30).
			Background(background).
			Border(lipgloss.InnerHalfBlockBorder(), false, false, false, true).
			BorderForeground(catppuccingo.Mocha.Mauve()).
			PaddingLeft(1)
	}
	bullet := lipgloss.NewStyle().
		Foreground(colors[task.status]).
		Bold(true)
	bars := lipgloss.NewStyle().
		Foreground(catppuccingo.Mocha.Surface2())
	text := lipgloss.NewStyle()

	if highlight {
		bullet = bullet.Background(background)
		bars = bars.Background(background)
		text = text.Background(background)
	}

	str := lipgloss.JoinHorizontal(
		lipgloss.Left,
		bars.Render("["),
		bullet.Render(bullets[task.status]),
		bars.Render("]"),
		text.Render(" "),
		text.Render(task.imperative),
	)

	fmt.Fprint(w, box.Render(str))
}

var bullets = map[status]string{
	Incomplete: " ",
	InProgress: "/",
	Cancelled:  "-",
	Complete:   "x",
}

var colors = map[status]color.Color{
	Incomplete: lipgloss.BrightBlue,
	InProgress: lipgloss.BrightYellow,
	Cancelled:  lipgloss.BrightRed,
	Complete:   lipgloss.BrightGreen,
}

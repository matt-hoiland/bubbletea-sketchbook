package window

import (
	"fmt"
	"image/color"
	"io"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	catppuccingo "github.com/catppuccin/go"
)

type taskDelegateKeyMap struct {
	CompleteTask key.Binding
	BeginTask    key.Binding
	CancelTask   key.Binding
	UndoTask     key.Binding
	MoveTaskUp   key.Binding
	MoveTaskDown key.Binding
	DeleteTask   key.Binding
	NewTask      key.Binding
}

var (
	bullets = map[status]string{
		Incomplete: " ",
		InProgress: "/",
		Cancelled:  "-",
		Complete:   "x",
	}

	colors = map[status]color.Color{
		Incomplete: lipgloss.BrightBlue,
		InProgress: lipgloss.BrightYellow,
		Cancelled:  lipgloss.BrightRed,
		Complete:   lipgloss.BrightGreen,
	}

	taskDelegateKeys = taskDelegateKeyMap{
		CompleteTask: key.NewBinding(
			key.WithKeys("space"),
			key.WithHelp("space", "complete task"),
		),
		BeginTask: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "begin task"),
		),
		CancelTask: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "cancel task"),
		),
		UndoTask: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "undo task"),
		),
		MoveTaskUp: key.NewBinding(
			key.WithKeys("shift+up", "K"),
			key.WithHelp("shift+↑/k", "move task up"),
		),
		MoveTaskDown: key.NewBinding(
			key.WithKeys("shift+down", "J"),
			key.WithHelp("shift+↓/j", "move task down"),
		),
		DeleteTask: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete task"),
		),
		NewTask: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new task"),
		),
	}
)

type TaskDelegate struct {
	keys taskDelegateKeyMap
}

func NewTaskDelegate() TaskDelegate {
	return TaskDelegate{
		keys: taskDelegateKeys,
	}
}

// FullHelp implements [help.KeyMap].
func (t TaskDelegate) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			t.keys.CompleteTask,
			t.keys.UndoTask,
			t.keys.BeginTask,
			t.keys.CancelTask,
		},
		{
			t.keys.MoveTaskUp,
			t.keys.MoveTaskDown,
			t.keys.NewTask,
			t.keys.DeleteTask,
		},
	}
}

// ShortHelp implements [help.KeyMap].
func (t TaskDelegate) ShortHelp() []key.Binding {
	return []key.Binding{
		t.keys.CompleteTask,
		t.keys.UndoTask,
	}
}

func (TaskDelegate) Height() int {
	return 1
}

func (TaskDelegate) Spacing() int {
	return 0
}

func (t TaskDelegate) Update(msg tea.Msg, ll *list.Model) tea.Cmd {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, t.keys.CompleteTask):
			cmds = append(cmds, t.setItemStatus(ll, Complete))
		case key.Matches(msg, t.keys.BeginTask):
			cmds = append(cmds, t.setItemStatus(ll, InProgress))
		case key.Matches(msg, t.keys.CancelTask):
			cmds = append(cmds, t.setItemStatus(ll, Cancelled))
		case key.Matches(msg, t.keys.UndoTask):
			cmds = append(cmds, t.setItemStatus(ll, Incomplete))
		case key.Matches(msg, t.keys.MoveTaskUp):
			cmds = append(cmds, t.moveTask(ll, -1))
		case key.Matches(msg, t.keys.MoveTaskDown):
			cmds = append(cmds, t.moveTask(ll, +1))
		case key.Matches(msg, t.keys.DeleteTask):
			cmds = append(cmds, t.deleteTask(ll))
		case key.Matches(msg, t.keys.NewTask):
			cmds = append(cmds, t.newTask(ll))
		}
	}
	return tea.Batch(cmds...)
}

func (TaskDelegate) deleteTask(ll *list.Model) tea.Cmd {
	ll.RemoveItem(ll.GlobalIndex())
	return nil
}

func (TaskDelegate) moveTask(ll *list.Model, dir int) tea.Cmd {
	i := ll.GlobalIndex()
	if i+dir < 0 || i+dir >= len(ll.Items()) {
		return nil
	}
	n := i + dir
	item := ll.SelectedItem()
	ll.RemoveItem(i)
	cmd := ll.InsertItem(n, item)
	ll.Select(n)
	return cmd
}

func (TaskDelegate) newTask(ll *list.Model) tea.Cmd {
	cmd := ll.InsertItem(ll.GlobalIndex()+1, Task{status: Incomplete, imperative: "Define this new task"})
	ll.Select(ll.GlobalIndex() + 1)
	return cmd
}

func (TaskDelegate) setItemStatus(ll *list.Model, s status) tea.Cmd {
	i := ll.GlobalIndex()
	item, ok := ll.SelectedItem().(Task)
	before := item
	if !ok {
		return nil
	}
	item.status = s
	after := item
	return tea.Batch(
		ll.SetItem(i, item),
		taskChanged(before, after),
	)
}

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

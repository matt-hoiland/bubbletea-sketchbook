package window

import (
	"fmt"
	"image/color"
	"io"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
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
	EditTask     key.Binding
	ConfirmEdit  key.Binding
	CancelEdit   key.Binding
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
		EditTask: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit task"),
		),
		ConfirmEdit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
			key.WithDisabled(),
		),
		CancelEdit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
			key.WithDisabled(),
		),
	}
)

type TaskDelegate struct {
	keys taskDelegateKeyMap

	editor       textinput.Model
	editorActive bool
}

func NewTaskDelegate() *TaskDelegate {
	editor := textinput.New()
	editor.Prompt = ""

	return &TaskDelegate{
		keys:         taskDelegateKeys,
		editor:       editor,
		editorActive: false,
	}
}

// FullHelp implements [help.KeyMap].
func (t *TaskDelegate) FullHelp() [][]key.Binding {
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
			t.keys.EditTask,
			t.keys.ConfirmEdit,
			t.keys.CancelEdit,
		},
	}
}

// ShortHelp implements [help.KeyMap].
func (t *TaskDelegate) ShortHelp() []key.Binding {
	return []key.Binding{
		t.keys.CompleteTask,
		t.keys.UndoTask,
	}
}

func (*TaskDelegate) Height() int {
	return 1
}

func (*TaskDelegate) Spacing() int {
	return 0
}

func (t *TaskDelegate) Update(msg tea.Msg, ll *list.Model) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, t.keys.CompleteTask):
			return t.setItemStatus(ll, Complete)
		case key.Matches(msg, t.keys.BeginTask):
			return t.setItemStatus(ll, InProgress)
		case key.Matches(msg, t.keys.CancelTask):
			return t.setItemStatus(ll, Cancelled)
		case key.Matches(msg, t.keys.UndoTask):
			return t.setItemStatus(ll, Incomplete)
		case key.Matches(msg, t.keys.MoveTaskUp):
			return t.moveTask(ll, -1)
		case key.Matches(msg, t.keys.MoveTaskDown):
			return t.moveTask(ll, +1)
		case key.Matches(msg, t.keys.DeleteTask):
			return t.deleteTask(ll)
		case key.Matches(msg, t.keys.NewTask):
			return t.newTask(ll)
		case key.Matches(msg, t.keys.EditTask):
			return t.editTask(ll)
		case key.Matches(msg, t.keys.ConfirmEdit):
			return t.confirmEdit(ll)
		case key.Matches(msg, t.keys.CancelEdit):
			return t.cancelEdit()
		}
	}

	var cmd tea.Cmd
	t.editor, cmd = t.editor.Update(msg)

	return cmd
}

func (*TaskDelegate) deleteTask(ll *list.Model) tea.Cmd {
	ll.RemoveItem(ll.GlobalIndex())
	return nil
}

func (t *TaskDelegate) cancelEdit() tea.Cmd {
	t.editor.Blur()
	t.activateEditor(false)
	return nil
}

func (t *TaskDelegate) confirmEdit(ll *list.Model) tea.Cmd {
	defer t.cancelEdit()
	item, ok := ll.SelectedItem().(Task)
	if !ok {
		return nil
	}
	item.imperative = t.editor.Value()
	return ll.SetItem(ll.GlobalIndex(), item)
}

func (t *TaskDelegate) editTask(ll *list.Model) tea.Cmd {
	item, ok := ll.SelectedItem().(Task)
	if !ok {
		return nil
	}

	t.activateEditor(true)

	t.editor.SetValue(item.imperative)
	t.editor.CursorEnd()
	return t.editor.Focus()
}

func (t *TaskDelegate) activateEditor(active bool) {
	t.editorActive = active

	// disable many controls
	t.keys.CompleteTask.SetEnabled(!active)
	t.keys.UndoTask.SetEnabled(!active)
	t.keys.BeginTask.SetEnabled(!active)
	t.keys.CancelTask.SetEnabled(!active)
	t.keys.MoveTaskDown.SetEnabled(!active)
	t.keys.MoveTaskUp.SetEnabled(!active)
	t.keys.EditTask.SetEnabled(!active)
	t.keys.NewTask.SetEnabled(!active)

	// enable others
	t.keys.ConfirmEdit.SetEnabled(active)
	t.keys.CancelEdit.SetEnabled(active)
}

func (*TaskDelegate) moveTask(ll *list.Model, dir int) tea.Cmd {
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

func (*TaskDelegate) newTask(ll *list.Model) tea.Cmd {
	cmd := ll.InsertItem(ll.GlobalIndex()+1, Task{status: Incomplete, imperative: "Define this new task"})
	ll.Select(ll.GlobalIndex() + 1)
	return cmd
}

func (*TaskDelegate) setItemStatus(ll *list.Model, s status) tea.Cmd {
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

func (t *TaskDelegate) Render(w io.Writer, ll list.Model, index int, item list.Item) {
	task, ok := item.(Task)
	if !ok {
		task = Task{status: Incomplete, imperative: "unknown task"}
	}

	var (
		background = catppuccingo.Mocha.Surface0()
		selector   = catppuccingo.Mocha.Mauve()

		imperative = task.imperative
	)

	if t.editorActive && index == ll.GlobalIndex() {
		selector = catppuccingo.Mocha.Peach()
		imperative = t.editor.View()
	}

	var (
		highlight = false
		box       = lipgloss.NewStyle().PaddingLeft(2)
	)
	if index == ll.Index() {
		highlight = true
		box = lipgloss.NewStyle().
			Width(30).
			Background(background).
			Border(lipgloss.InnerHalfBlockBorder(), false, false, false, true).
			BorderForeground(selector).
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
		text.Render(imperative),
	)

	fmt.Fprint(w, box.Render(str))
}

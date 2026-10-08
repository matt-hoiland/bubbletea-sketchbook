package window

import tea "charm.land/bubbletea/v2"

type status int

const (
	Incomplete status = iota
	InProgress
	Complete
	Cancelled
)

func (s status) String() string {
	switch s {
	case Incomplete:
		return "undone"
	case InProgress:
		return "started"
	case Complete:
		return "completed"
	case Cancelled:
		return "cancelled"
	}
	return "unknown status"
}

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

type TaskChangedMsg struct {
	Before, After Task
}

func taskChanged(before, after Task) tea.Cmd {
	return func() tea.Msg {
		return TaskChangedMsg{
			Before: before,
			After:  after,
		}
	}
}

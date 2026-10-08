package window

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/paginator"
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
	}, NewTaskDelegate(), 0, 0)

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

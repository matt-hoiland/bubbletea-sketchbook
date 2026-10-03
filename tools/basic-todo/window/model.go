package window

import (
	"charm.land/bubbles/v2/help"
)

type Model struct {
	keys keyMap
	help help.Model

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

	m := &Model{
		keys: keys,
		help: help.New(),

		title: c.title,
	}

	return m
}

func WithTitle(title string) Option {
	return func(c *conf) {
		c.title = title
	}
}

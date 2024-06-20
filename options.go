package main

type Option struct {
	blame bool
}

func NewOption() *Option {
	return &Option{
		blame: false,
	}
}

package main

type Option struct {
	blame    bool
	remote   bool
	showHelp bool
}

func NewOption() *Option {
	return &Option{
		blame:    false,
		remote:   false,
		showHelp: false,
	}
}
